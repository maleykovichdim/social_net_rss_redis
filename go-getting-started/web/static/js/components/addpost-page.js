import { doGet, doPost } from "../http.js"
import renderList from "./list.js";
import renderUser from "./user-profile.js";
import { navigate } from "../lib/router.js";

const PAGE_SIZE = 3
// var search_mode = "Surname"
const template = document.createElement("template")
template.innerHTML = `
    <div class="container">

        <h1>Test Social Net - Add Post </h1>
        <h1> </h1>
        <p><div>
        <textarea class="textarea2" id="post" name="post" placeholder="your post" rows="10"  cols="50"   ></textarea>
        <p><button id="add_post" >Save post</button></p>
    </div>
`

export default async function renderAddPostPage() {
    const url = new URL(location.toString())
    const page = /** @type {DocumentFragment} */ (template.content.cloneNode(true))
    const button = /** @type {HTMLButtonElement} */ (page.getElementById("add_post"))

    page.getElementById("post").style.height = "100%";
    button.addEventListener("click", onSavePostSubmit)
    return page
}

 /**
 * @param {Event} ev
 */
  async function onSavePostSubmit(ev) {
    ev.preventDefault() 
    const input = /** @type {HTMLTextAreaElement} */ (document.getElementById("post"))
    var text = input.value
    var res = await savePost(text)
    console.log(res)
    if ( res.status == "Done"){
        location.assign('/') 
    }
 }


 /**
 * @returns {Promise<import('../types.js').StatusResponse>}
 * @param {any} text
 */
  function savePost(text) {
    return doPost('/api/auth_user/post', { text })
}


