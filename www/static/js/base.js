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

function calcExrId(exerEl) {
  let id = 0;
  for (const el of exerEl.querySelectorAll("input")) {
    id += el.dataset.awnser.length;
  }
  return id;
}

function calcScore(inputCount, correctCount) {
  if (inputCount === 0) {
    return "";
  }
  return (
    correctCount +
    "/" +
    inputCount +
    " " +
    Math.floor((correctCount / inputCount) * 100) +
    "%"
  );
}

function updateScore(exerciseEl, input = true) {
  const inputCount = exerciseEl.querySelectorAll("input").length;
  const correctCount =
    exerciseEl.querySelectorAll("input[data-correct=true]").length +
    exerciseEl.querySelectorAll(".correct").length;

  const score = (correctCount / inputCount) * 100;
  const starEl = exerciseEl.querySelector(".exr-footer .star");
  if (score == 100) {
    starEl.classList.add("fully-correct");
  } else {
    starEl.classList.remove("fully-correct");
  }

  exerciseEl.querySelector(".exr-footer .score").innerText =
    correctCount + "/" + inputCount;
  exerciseEl.querySelector(".exr-footer .score-bar .fill").style.width =
    score + "%";

  if (input) {
    localStorage.setItem(calcExrId(exerciseEl) + "_lastScore", score);
  }
  updateLastScore(exerciseEl);
}

function updateLastScore(exerciseEl) {
  const score = localStorage.getItem(calcExrId(exerciseEl) + "_lastScore");
  if (score) {
    exerciseEl.querySelector(".score").innerText =
      "laatste score " + Math.ceil(score) + "%";
  }
}

function fractionToNum(frac) {
  const slashIndex = frac.indexOf("/");
  if (slashIndex != -1) {
    const t = Number(frac.slice(0, slashIndex));
    const n = Number(frac.slice(slashIndex + 1));
    return t / n;
  }
  return Number(frac);
}

function checkNumInput(value, awnser) {
  value = value.replaceAll(" ", "");
  awnser = awnser.replaceAll(" ", "");
  if (value == awnser) {
    return true;
  }

  value = value.replace(",", ".");
  awnser = awnser.replace(",", ".");
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

$(".exr").on("input", function (_) {
  check(this);
});
$(".exr").each(function (_) {
  check(this);
});

function checkMultipleChoice(exrEl) {
  exrEl.classList.add("show-awnsers");
  for (const el of exrEl.querySelectorAll("li")) {
    const inputEl = el.querySelector("input");
    inputEl.dataset.correct =
      inputEl.dataset.awnser === String(inputEl.checked);
  }
}

$(".exr-multiplechoice button").on("click", function (_) {
  checkMultipleChoice(this.parentElement);
  updateScore(this.parentElement.parentElement);
});

$(".exercise").on("input", function (_) {
  updateScore(this);
});
$(".exercise").each(function (_) {
  updateScore(this, false);
});

async function updateLessonPreview() {
  let lessonPreviewElement = document.getElementById("lessonPreview");
  let lessonMDElement = document.getElementById("lessonMD");
  let MDContent = lessonMDElement.innerText;

  let lessonPreviewHTML = await fetch("/lessons/preview", {
    method: "POST",
    headers: {
      "Content-Type": "text/plain",
    },
    body: MDContent,
  });

  lessonPreviewElement.innerHTML = lessonPreviewHTML;
}

async function saveLesson(id) {
  await fetch("/lessons/"+id, {
    method: "PUT",
    body: document.getElementById("LessonSrcCode").innerText,
  })
}
