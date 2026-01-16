function getSourceCode() {
  return document.getElementById("LessonSrcCode").innerText;
}

let previewedSrcCode=getSourceCode();
let dirty=false;

async function updateLessonPreview() {
  const srcCode = getSourceCode();
  if (previewedSrcCode == srcCode) return;
  const res = await fetch("/lessons/preview", {
    method: "POST",
    headers: {
      "Content-Type": "text/plain",
    },
    body: srcCode,
  });

  const htmlCode = await res.text();
  document.getElementById("lessonContent").innerHTML = htmlCode;
  previewedSrcCode = srcCode;
}

setInterval(async () => {
  await updateLessonPreview(); 
}, 1000);

async function saveLesson() {
  setDirty(false);
  await fetch("/lessons/" + LESSON_ID, {
    method: "PUT",
    body: getSourceCode(),
  });
}

document.getElementById("LessonSrcCode").addEventListener("input", (e) => {
  setDirty(true);
});

function setDirty(value) {
  if (dirty === value) { return; }
  dirty=value;
  document.getElementById("unsaved-icon").style.display=dirty?"inline":"none";
  if (value){
    window.addEventListener("beforeunload", onBeforeUnload)
  } else {
    window.removeEventListener("beforeunload", onBeforeUnload)
  }
}

function onBeforeUnload(e) {
  e.preventDefault();
}

document.addEventListener("keydown", async (e) => {
  if (e.key == "s" && e.ctrlKey) {
    e.preventDefault();
    await saveLesson();
  }
})

