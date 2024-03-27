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
    event.preventDefault(); // prevent default form submission
    var author = document.getElementById("author").value;
    var siteTitle = document.getElementById("siteTitle").value;
    var baseURL = document.getElementById("baseURL").value;
    if (!author || !siteTitle || !baseURL) {
        alert("Please fill out all fields.");
        return;
    }
    var formData = {
        "author": author,
        "siteTitle": siteTitle,
        "baseURL": baseURL
    };
    var jsonData = JSON.stringify(formData);
    localStorage.setItem("formData", jsonData);
    showSlide(slides.length - 1);
    var today = new Date();
    var dateString = today.getFullYear() + '-' + (today.getMonth() + 1) + '-' + today.getDate();
    var fileName = 'site-' + dateString + '.json';
    var blob = new Blob([jsonData], { type: 'application/json' });
    var url = URL.createObjectURL(blob);
    var a = document.createElement('a');
    a.href = url;
    a.download = fileName;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
}
showSlide(0); // show the first slide