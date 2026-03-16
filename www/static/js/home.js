class LessonCategory {
  element;
  categoryId;
  isOpen;
  constructor(categoryId, element) {
    this.categoryId = categoryId;
    this.element = element;
    this.isOpen = true;
    this.initialiseEventListeners();
  }
  close() {
    this.element.classList.add("hidden");
  }

  open() {
    this.element.classList.remove("hidden");
  }

  toggle() {
    this.isOpen = !this.isOpen;
    this.isOpen ? this.open() : this.close();
  }

  initialiseEventListeners() {
    let collapseButton = this.element.querySelector(".collapse-button");
    if (!collapseButton) return;
    collapseButton.addEventListener("click", () => this.toggle());
  }
}

class Lesson {
  element;
  lessonId;
  constructor(lessonId, element) {
    this.lessonId = lessonId;
    this.element = element;
    this.initialiseEventListeners();
  }

  async saveLocation() {}

  drop() {}

  move(x, y) {
    this.element.style.left = x + "px";
    this.element.style.top = y + "px";
  }

  startMoving() {
    this.element.classList.add("is-moving");
    this.element.parentElement = document.body;
    let moveEventHandler = (e) => {
      this.move(e.clientX, e.clientY);
    };
    document.addEventListener("mousemove", moveEventHandler);
    document.addEventListener("mouseup", () => {
      document.removeEventListener("mousemove", moveEventHandler);
      this.drop();
    });
  }

  initialiseEventListeners() {
    let moveButton = this.element.querySelector(".move-button");
    if (!moveButton) return;

    moveButton.addEventListener("click", () => this.startMoving());
  }
}

let lessonCategories = createCategoryObjects();
function createCategoryObjects() {
  let lessons = Array.from(document.querySelectorAll(".lesson-category")).map(
    (lessonCategory) => {
      let lessonId = lessonCategory.querySelector(".chapter").id.split("-")[1];
      return new LessonCategory(lessonId, lessonCategory);
    },
  );
  return lessons;
}
