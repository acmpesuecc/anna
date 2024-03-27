let currentSlide = 0;
const slides = document.querySelectorAll('.content');

function showSlide(index) {
    slides.forEach((slide, i) => {
        if (i === index) {
            slide.style.display = 'block';
        } else {
            slide.style.display = 'none';
        }
    });
}

function nextSlide() {
    currentSlide++;
    if (currentSlide >= slides.length) {
        currentSlide = slides.length - 1;
    }
    showSlide(currentSlide);
}

function prevSlide() {
    currentSlide--;
    if (currentSlide < 0) {
        currentSlide = 0;
    }
    showSlide(currentSlide);
}

function submitForm() {
    var author = document.getElementById("author").value;
    var siteTitle = document.getElementById("siteTitle").value;
    var baseURL = document.getElementById("baseURL").value;

    if (!author || !siteTitle || !baseURL) {
        alert("Please fill out all fields.");
        return;
    }

    var formData = '{' +
        `"author":"${author}",` +
        `"siteTitle":"${siteTitle}",` +
        `"baseURL":"${baseURL}"` +
        '}';

    nextSlide(); // Move to the next slide after form validation

    fetch('/submit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: formData
    })

    setTimeout(() => {
        window.location.href = 'http://localhost:8000';
    }, 3000); // 3s
}

showSlide(0);
