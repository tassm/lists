const loadHtml = async (path) => {
  const response = await fetch(path);
  const text = await response.text();
  return text;
}

document.addEventListener('DOMContentLoaded', async () => {
  var content = document.getElementById('content');

  // load the list content
  var listPage = await loadHtml("list.html")
  var content = document.getElementById('content');
  content.innerHTML = listPage;

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