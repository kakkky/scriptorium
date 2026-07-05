import "https://esm.sh/@hotwired/turbo@8.0.4"
import { Application, Controller } from "https://esm.sh/@hotwired/stimulus@3.2.2"

class ClearInputController extends Controller {
  static targets = ["input"]

  clear() {
    this.inputTarget.value = ""
    this.inputTarget.focus()
  }
}

const app = Application.start()
app.register("clear-input", ClearInputController)

document.addEventListener("turbo:submit-end", (e) => {
  if (e.detail.success && e.target.matches("[data-reset-on-submit]")) {
    e.target.reset()
  }
})
