const DEBUG = true;

if (DEBUG) {
  const cssURL = "/static/css/base.css"; // change to your CSS file path
  const refreshRate = 500; // milliseconds

  function reloadCSS() {
    const linkId = "live-css-reload";
    let oldLink = document.getElementById(linkId);
    const newLink = document.createElement("link");
    newLink.rel = "stylesheet";
    newLink.id = linkId;
    // Append a timestamp to bypass cache
    newLink.href = cssURL + "?t=" + Date.now();

    newLink.onload = () => {
      if (oldLink) oldLink.remove();
    };

    document.head.appendChild(newLink);
  }

  reloadCSS();
  setInterval(reloadCSS, refreshRate);
}

function fractionToNum(frac) {
  const slashIndex = frac.indexOf("/");
  if (slashIndex != -1) {
    const t = Number(frac.slice(0, slashIndex))
    const n = Number(frac.slice(slashIndex + 1))
    return t / n
  }
  return Number(frac)
}

function checkNumInput(value, awnser) {
  value = value.replaceAll(" ", "");
  awnser = awnser.replaceAll(" ", "");
  if (value == awnser) {
    return true;
  }

  value = value.replace(",", ".")
  awnser = awnser.replace(",", ".")
  if (value == awnser) {
    return true;
  }
  if (fractionToNum(awnser) == fractionToNum(value)) {
    return true;
  }
  return false;
}

function check(inputEl) {
  let value = inputEl.value;

  let correct = value == inputEl.dataset.awnser;
  if (inputEl.dataset.type == "num") {
    correct = checkNumInput(value, inputEl.dataset.awnser);
  }

  if (correct) {
    inputEl.classList.add("correct");
  } else {
    inputEl.classList.remove("correct");
  }
}


$(".exr").on("input", function(_) {
  check(this)
})
$(".exr").each(function(_) {
  check(this)
})
