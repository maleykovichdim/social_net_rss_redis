import { doGet, doPost } from "../http.js"
// import renderList from "./list.js";
// import renderUser from "./user-profile.js";
// import { navigate } from "../lib/router.js";

const PAGE_SIZE = 3
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
    const page = /** @type {DocumentFragment} */ (template.content.cloneNode(true))
    var posts = await getRssPosts()
    var d1 = /** @type {HTMLDivElement} */page.getElementById('output');
    for (var post of posts) {
        console.log(post)
        d1.insertAdjacentHTML('beforebegin', 
        '<p><div class="div-3">Post: '+ post.id +' Content:</div><div class="div-1">'+post.content+'</div>'+
        '<div class="div-2">'+post.created_at+'  author:'+post.author_id);  
    }
    return page
}


/**
 * @returns {Promise <import('../types.js').Post[]>}
 */
 async function getRssPosts() {
    return doGet('/api/auth_user/rss_feed')
}


