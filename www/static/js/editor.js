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
  await fetch("/lesson/" + id, {
    method: "PUT",
    body: document.getElementById("lessonSrcCode").innerText,
  });
}

async function loadContent(lessonId) {
  let lessonMD = await fetch(`/lessons/${lessonId}/src`, {
    method: "GET",
  });
  document.getElementById("content").innerText = lessonMD;
}
