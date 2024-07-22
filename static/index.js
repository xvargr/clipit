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
    .then(async (res) => {
      if (res.status !== 200) throw res;
      return await res.json();
    })
    .then((data) => setToReady(data))
    .catch(async (err) => {
      window.alert(await err.text());
    });
});

UrlInput.addEventListener("input", setToStandby);

copyButton.addEventListener("click", copyToClipboard);

function setToStandby() {
  UrlDisplay.innerText = "...";
  copyButton.disabled = true;
}

function setToReady({ shortenedUrl, validity } = dat) {
  UrlDisplay.innerText = shortenedUrl;
  UrlDisplay.href = shortenedUrl;
  copyButton.disabled = false;
  // console.log(parseValidityMessage(validity));
}

// validity is hardcoded for now
// function parseValidityMessage(rawString) {
//   // raw string is in the format of <number>h<number>m<number>s
//   const [hours, minutes, seconds] = rawString.split(/[hms]/).map(Number);

//   return (
//     "This link will expire in " +
//     (hours ? `${hours} hour${hours > 1 ? "s" : ""} ` : "") +
//     (minutes ? `${minutes} minute${minutes > 1 ? "s" : ""} ` : "") +
//     (seconds ? `${seconds} second${seconds > 1 ? "s" : ""}` : "")
//   );
// }

function copyToClipboard() {
  navigator.clipboard?.writeText(UrlDisplay.innerText);
}

(function buttonClickBounce() {
  const buttons = document.querySelectorAll("button");

  buttons.forEach((button) => {
    button.addEventListener("mousedown", () => {
      button.classList.add("clicking");
    });

    button.addEventListener("mouseup", () => {
      button.classList.remove("clicking");
    });
  });
})();
