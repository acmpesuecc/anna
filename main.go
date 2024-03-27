package main

import (
	"html/template"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/acmpesuecc/anna/cmd/ssg"
	"github.com/spf13/cobra"
)

func main() {
	var serve bool
	var webconsole bool
	var addr string
	var draft bool
	var validateHTML bool
	var prof bool
	StartProfiling()
	startTime := time.Now()
	rootCmd := &cobra.Command{
		Use:   "anna",
		Short: "Static Site Generator",
		Run: func(cmd *cobra.Command, args []string) {

			generator := ssg.Generator{
				ErrorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
				Templates:   make(map[template.URL]ssg.TemplateData),
				TagsMap:     make(map[string][]ssg.TemplateData),
			}

			if draft {
				generator.RenderDrafts = true
			}

			if serve {
				if prof {
					go func() {
						for {
							time.Sleep(5 * time.Second) //change as per needed
							PrintStats(time.Since(startTime))
						}
					}()
				}
				generator.StartLiveReload(addr)
			}

			if !prof {
				generator.RenderSite("")
			}

			if validateHTML {
				ssg.ValidateHTMLContent()
			}
			if prof {

				generator.RenderSite("")

				elapsedTimesince := time.Since(startTime) //this didn't work for some reason and was giving negitive deviation
				// elapsedTime := time.Now().Sub(startTime)

				PrintStats(elapsedTimesince)
				// PrintStats(elapsedTime)
				defer StopProfiling()
			}
			if webconsole {
				ssg.Wizard()
			}

		},
	}

	rootCmd.Flags().BoolVarP(&serve, "serve", "s", false, "serve the rendered content")
	rootCmd.Flags().BoolVarP(&webconsole, "webconsole", "w", false, "wizard to setup anna")
	rootCmd.Flags().StringVarP(&addr, "addr", "a", "8000", "ip address to serve rendered content to")
	rootCmd.Flags().BoolVarP(&draft, "draft", "d", false, "renders draft posts")
	rootCmd.Flags().BoolVarP(&validateHTML, "validate-html", "v", false, "validate semantic HTML")
	rootCmd.Flags().BoolVarP(&prof, "prof", "p", false, "enable profiling")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}

func StartProfiling() {
	runtime.MemProfileRate = 1
	go func() {
		for {
			time.Sleep(5 * time.Second)
			memProfile := pprof.Lookup("heap")
			if memProfile != nil {
				file, err := os.Create("mem_profile.pprof")
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
				memProfile.WriteTo(file, 1)
			}
			cpuProfile := pprof.Lookup("goroutine")
			if cpuProfile != nil {
				file, err := os.Create("cpu_profile.pprof")
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
				cpuProfile.WriteTo(file, 1)
			}
		}
	}()
}

func PrintStats(elapsedTime time.Duration) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	log.Printf("Memory Usage: %d bytes", memStats.Alloc)
	log.Printf("Time Elapsed: %s", elapsedTime)
	cpuUsage := runtime.NumCPU()
	threads := runtime.GOMAXPROCS(0)
	// memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	log.Printf("Threads: %d", threads)
	log.Printf("Cores: %d", cpuUsage)
	log.Printf("Time Taken: %s", elapsedTime)
	log.Printf("Allocated Memory: %d bytes", memStats.Alloc)
	log.Printf("Total Memory Allocated: %d bytes", memStats.TotalAlloc)
	log.Printf("Heap Memory In Use: %d bytes", memStats.HeapInuse)
	log.Printf("Heap Memory Idle: %d bytes", memStats.HeapIdle)
	log.Printf("Heap Memory Released: %d bytes", memStats.HeapReleased)
	log.Printf("Number of Goroutines: %d", runtime.NumGoroutine())
	print("-----------------------------------\n")
}

func StopProfiling() {
	pprof.StopCPUProfile()
}
