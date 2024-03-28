let currentSlide = 0;
const slides = document.querySelectorAll('.content');

function showSlide(index) {
    slides.forEach((slide, i) => {
        slide.style.display = i === index ? 'block' : 'none';
    });
}

function nextSlide() {
    currentSlide = Math.min(currentSlide + 1, slides.length - 1);
    showSlide(currentSlide);
}

function prevSlide() {
    currentSlide = Math.max(currentSlide - 1, 0);
    showSlide(currentSlide);
}

function checkFormValidity() {
    var author = document.getElementById("author").value;
    var siteTitle = document.getElementById("siteTitle").value;
    var baseURL = document.getElementById("baseURL").value;
    var themeURL = document.getElementById("themeURL").value;

    var nextButton = document.getElementById("nextButton");
    var nextButton2 = document.getElementById("nextButton2");
    var nextButton3 = document.getElementById("nextButton3");

    nextButton.disabled = !(author);
    nextButton2.disabled = !(siteTitle);
    nextButton3.disabled = !(baseURL);
}

function submitForm() {
    var author = document.getElementById("author").value;
    var siteTitle = document.getElementById("siteTitle").value;
    var baseURL = document.getElementById("baseURL").value;
    var themeURL = document.getElementById("themeURL").value;

    if (!author.trim() || !siteTitle.trim() || !baseURL.trim() || !themeURL.trim()) {
        alert("Please fill out all fields.");
        return;
    }

    var formData = JSON.stringify({
        "author": author,
        "siteTitle": siteTitle,
        "baseURL": baseURL,
        "themeURL": themeURL,
        "navbar": "index,about"
    });

    showSlide(slides.length - 1);
    fetch('/submit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: formData
    });

    setTimeout(() => {
        window.location.href = 'http://localhost:8000';
    }, 3000); // 3s
}

showSlide(0);
