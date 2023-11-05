const loadAndCacheHtmlFile = async (path) => {
  const response = await fetch(path);
  const html = await response.text();
  const parser = new DOMParser();
  return parser.parseFromString(html, 'text/html')
}

document.addEventListener('DOMContentLoaded', async () => {
  var content = document.getElementById('content');

  // load the list content
  var listDom = await loadAndCacheHtmlFile('list.html');
  var content = document.getElementById('content');
  content.innerHTML = listDom.getElementById('list-container').outerHTML;

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
    }
]

  listData.forEach(element => {
    console.log(element['item'])
    var listCard = listDom.getElementById('list-item');
    listCard.getElementById("item").innerText = element['item'];
    document.getElementById('list-container').appendChild(listCard)
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