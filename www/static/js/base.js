const DEBUG = false;

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

function updateScore(exerciseEl) {
  const inputCount = exerciseEl.querySelectorAll("input").length;
  const correctCount = exerciseEl.querySelectorAll(".correct").length;

  console.log(inputCount, correctCount);
  exerciseEl.querySelector(".score").innerText = correctCount + "/" + inputCount + " " + Math.floor((correctCount / inputCount) * 100) + "%";
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

  if (value.trim() == "") {
    inputEl.classList.remove("correct");
    inputEl.classList.remove("wrong");
    return;
  }

  let correct = value == inputEl.dataset.awnser;
  if (inputEl.dataset.type == "num") {
    correct = checkNumInput(value, inputEl.dataset.awnser);
  }

  if (correct) {
    inputEl.classList.add("correct");
    inputEl.classList.remove("wrong");
  } else {
    inputEl.classList.add("wrong");
    inputEl.classList.remove("correct");
  }
}

$(".exr").on("input", function(_) {
  check(this)
})
$(".exr").each(function(_) {
  check(this)
})

function checkMultipleChoice(exrEl) {
  for (const el of exrEl.querySelectorAll("li")) {
    const inputEl = el.querySelector("input");
    if (String(inputEl.checked) == inputEl.dataset.awnser) {
      el.classList.add("correct");
      el.classList.remove("wrong");
    } else {
      el.classList.add("wrong");
      el.classList.remove("correct");
    }
  }
}

$(".exr-multiplechoice button").on("click", function(_) {
  checkMultipleChoice(this.parentElement)
  updateScore(this.parentElement.parentElement);
})

$(".exercise").on("input", function(_) {
  updateScore(this);
})
$(".exercise").each(function(_) {
  updateScore(this);
})
