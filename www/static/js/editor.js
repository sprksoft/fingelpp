let exerciseSvg = `<svg width="18px" style="scale:1.2" height="18px" viewBox="0 0 24 24" fill="none">
<path d="M8 12.3333L10.4615 15L16 9M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z"  stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
</svg>`;
let infoSvg = `<svg width="18px" style="scale:1.2" height="18px" viewBox="0 0 24 24" fill="none">
<path d="M12 8H12.01M12 11V16M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
</svg>`;

var fingelEditor = new TinyMDE.Editor({
  element: "LessonSrcCode",
  content: ORIGINAL_SOURCE,
});
var commandBar = new TinyMDE.CommandBar({
  element: "LessonToolbar",
  editor: fingelEditor,
  commands: [
    "bold",
    "italic",
    "|",
    {
      name: "exercise",
      title: "exercise",
      innerHTML: exerciseSvg,
      action: addExercise,
    },
    {
      name: "info",
      title: "info",
      innerHTML: infoSvg,
      action: addInfo,
    },
    "|",
    "h1",
    "h2",
    "|",
    "ul",
    "|",
    "insertLink",
    "|",
    "undo",
    "redo",
  ],
});

function addExercise() {
  fingelEditor.paste(`
> [EX] 
> 
>
> @[o] 
> @[x] 
> @[x] 
`);
}

function addInfo() {
  fingelEditor.paste(`
> [INFO] 
> 
>
`);
}

function getSourceCode() {
  return fingelEditor.getContent();
}

let previewedSrcCode = getSourceCode();
let dirty = false;

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
  updateExercises();
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
  if (dirty === value) {
    return;
  }
  dirty = value;
  document.getElementById("unsaved-icon").style.display = dirty
    ? "inline"
    : "none";
  if (value) {
    window.addEventListener("beforeunload", onBeforeUnload);
  } else {
    window.removeEventListener("beforeunload", onBeforeUnload);
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
});
