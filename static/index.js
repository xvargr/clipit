const form = document.querySelector("form");
const UrlInput = document.querySelector("#url-input");
const UrlDisplay = document.querySelector("#url");
const copyButton = document.querySelector("#copy-button");

form.addEventListener("submit", (e) => {
  e.preventDefault();

  const formData = new FormData(form);
  fetch("/shorten", {
    method: "POST",
    body: formData,
  })
    .then((res) => res.json())
    .then((data) => setToReady(data))
    .catch((err) => {
      console.error(err);
    });
});

UrlInput.addEventListener("input", setToStandby);

copyButton.addEventListener("click", copyToClipboard);

function setToStandby() {
  UrlDisplay.innerText = "...";
  copyButton.disabled = true;
}

function setToReady(shortenedUrl) {
  UrlDisplay.innerText = shortenedUrl;
  UrlDisplay.href = shortenedUrl;
  copyButton.disabled = false;
}

function copyToClipboard() {
  navigator.clipboard?.writeText(UrlDisplay.innerText);
}
