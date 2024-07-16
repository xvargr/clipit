const form = document.querySelector("form");
form.addEventListener("submit", (e) => {
  e.preventDefault();

  const formData = new FormData(form);
  console.log(formData.get("url"));

  fetch("/shorten", {
    method: "POST",
    body: formData,
  })
    .then((res) => res.json())
    .then((data) => {
      const url = document.querySelector("#url");
      url.innerText = data;
      url.href = data;
    })
    .catch((err) => {
      console.error(err);
    });
});
