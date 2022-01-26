import { doGet, doPost } from "../http.js"
import renderList from "./list.js";
import renderUser from "./user-profile.js";
import { navigate } from "../lib/router.js";

const PAGE_SIZE = 3
// var search_mode = "Surname"
const template = document.createElement("template")
template.innerHTML = `

    <div class="container">

        <h1>Test Social Net - RSS </h1>
        <h1> </h1>
        <p>
        <div id="out"  style=
            "background-color: GRAY; overflow-y: scroll; overflow-x: scroll;
            height:400px; width: 100%; border:1px solid ;">
            <div id="output">
            </div> 
         </div>

    </div>
`

export default async function renderRssPage() {
    // const url = new URL(location.toString())
    const page = /** @type {DocumentFragment} */ (template.content.cloneNode(true))
    // const button = /** @type {HTMLButtonElement} */ (page.getElementById("add_post"))

    // page.getElementById("post").style.height = "100%";
    var posts = await getRssPosts()
    var d1 = /** @type {HTMLDivElement} */page.getElementById('output');
    for (var post of posts) {
        console.log(post)
        // d1.insertAdjacentHTML('beforebegin', '<div>'+JSON.stringify(post)+'<p></div><div></div>');
        d1.insertAdjacentHTML('beforebegin', 
        '<p><div class="div-3">Post: '+ post.id +' Content:</div><div class="div-1">'+post.content+'</div>'+
        '<div class="div-2">'+post.created_at+'  author:'+post.author_id);  
    }

    // page.root.innerHTML =  +
    // '<content></content>';


    // button.addEventListener("click", onSavePostSubmit)
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
    // location.assign('/addpost')
    // var interests = document.getElementById("interests")
    // var about = document.getElementById("personal")
    // await updatePersonalPage(interests.value, about.value)
 }


/**
 * @returns {Promise <import('../types.js').Post[]>}
 */
 async function getRssPosts() {
    return doGet('/api/auth_user/rss_feed')
}



// /**
//  * @returns {Promise<import('../types.js').StatusResponse>}
//  * @param {any} text
//  */
//    function getPosts(text) {
//     return doPost('/api/auth_user/post', { text })
// }





 /**
 * @returns {Promise<import('../types.js').StatusResponse>}
 * @param {any} text
 */
  function savePost(text) {
    return doPost('/api/auth_user/post', { text })
}



// /**
//  * @param {string} search
//  * @param {string=} after
//  * @returns {Promise<import("../types.js").User[]>}
//  */
// function fetchUsers(search, after = "") {
//     return doGet(`/api/users?search=${search}&after=${after}&first=${PAGE_SIZE}`)
// }

// /**
//  * @param {string} search
//  * @param {string=} after
//  * @returns {Promise<import("../types.js").User[]>}
//  */
//  function fetchUsersByInterests(search, after = "") {
//     return doGet(`/api/users_by_interests?search=${search}&after=${after}&first=${PAGE_SIZE}`)
// }