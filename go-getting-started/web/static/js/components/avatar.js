/**
 * @param {import("../types.js").User} user
 */
export default function renderAvatarHTML(user, title = "") {
    
    //TODO uncomment after solving problem with file saving
    //return 
    //user.avatarURL !== null && user.avatarURL !== undefined
    //    ? `<img class="avatar"  src="${user.avatarURL}" alt="${user.name}'s avatar" title="${title}">`
    //    :
    return     `<span class="avatar" data-initial="${user.name[0]}" title="${title}"></span>`
}
