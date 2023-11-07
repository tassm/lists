const loadAndCacheHtmlFile = async (path) => {
  const response = await fetch(path);
  const html = await response.text();
  const parser = new DOMParser();
  return parser.parseFromString(html, 'text/html')
}

const saveChanges = () => {
  console.log("hi")
}

document.addEventListener('DOMContentLoaded', async () => {
  var content = document.getElementById('content');

  // load the list content
  var listDom = await loadAndCacheHtmlFile('list.html');
  var content = document.getElementById('content');
  content.innerHTML = listDom.getElementById('list-template').innerHTML;

  // load the dummy data file
  const listData = [
    {
        "item": "milk",
        "done": false
    },
    {
        "item": "bread",
        "done": false
    },
    {
        "item": "servo pie",
        "done": false
    },
    {
      "item": "even more pie",
      "done": false
    },
    {
      "item": "a 3 legged dog",
      "done": false
    }
]

  var listTemplate = listDom.getElementById("list-item-template");
  listData.forEach(element => {
    console.log(element.item)
    listTemplate.content.querySelector("span").innerText = element.item
    var node = document.importNode(listTemplate.content.querySelector("div"), true)
    document.getElementById('list-container').appendChild(node)
  });

  

  // navBar.addEventListener('click', (event) => {
  // // Prevent the default link behavior
  // event.preventDefault();

  // // Get the target element of the click event
  // const target = event.target;

  // // If the target element is a link, update the content div to display the corresponding page
  // if (target.tagName === 'A') {
  //   const page = target.textContent.toLowerCase();

  //   if (page === 'home') {
  //     content.replaceChild(homePage, aboutPage);
  //   } else if (page === 'about') {
  //     content.replaceChild(aboutPage, homePage);
  //   }
  // }
  // });
});