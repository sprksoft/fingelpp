function getSourceCode() {
  return document.getElementById("LessonSrcCode").innerText;
}

let previewedSrcCode=getSourceCode();

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

async function saveLesson(id) {
  await fetch("/lessons/" + id, {
    method: "PUT",
    body: getSourceCode(),
  });
}
