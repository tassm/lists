var listDom;
var basePath;
var resetStyle;
var fullPath;

const endpoint = '/api/list';

const showErrorPopup = (errorMessage) => {
  // Set the error message.
  document.getElementById('error-message').innerHTML = errorMessage;
  // Display the error popup.
  document.getElementById('error-modal').classList.remove('hidden');
}

const hideErrorPopup = () => {
  // Hide the error modal.
  document.getElementById('error-modal').classList.add('hidden');
}

const postItem = async (item, list) => {
  // Create a new Fetch request object.
  const request = new Request(fullPath, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      item: item,
      listId: list
    }),
  });

  // Send the request and await the response.
  const response = await fetch(request);
}

// provide a list of item ids to update to done
const putCompletedItems = async (idList) => {
  // Create a new Fetch request object.
  const request = new Request(fullPath, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      ids: idList,
    }),
  });

  // Send the request and await the response.
  return await fetch(request);
}

const getItems = async () => {
  // Create a new Fetch request object.
  const request = new Request(fullPath, {
    method: 'GET'
  });

  // Send the request and await the response.
  const response = await fetch(fullPath);
}

const loadAndCacheHtmlFile = async (path) => {
  const response = await fetch(path);
  const html = await response.text();
  const parser = new DOMParser();
  return parser.parseFromString(html, 'text/html')
}

const newItem = () => {
  if (document.getElementById('new-item-text') == null) {
    var listTemplate = listDom.getElementById("list-item-new-template");
    var node = document.importNode(listTemplate.content.querySelector("div"), true)
    var container = document.getElementById('list-container')
    container.insertBefore(node, container.children[1])
  }
}

const saveNew = async () => {
  var itemText = document.getElementById('new-item-text').value
  var res = await postItem(itemText, "test-list")
  // CHECK IF IT IS A 201 CREATED


  var container = document.getElementById('list-container')
  container.children[1].remove()
  var listTemplate = listDom.getElementById("list-item-template");
  listTemplate.content.querySelector("span").innerText = itemText
  listTemplate.content.querySelector("input").id = "test-id"
  var node = document.importNode(listTemplate.content.querySelector("div"), true)
  container.insertBefore(node, container.children[1])
}

const saveChanges = async () => {
  const elements = document.querySelectorAll('input[type="checkbox"]');
  const completedIds = [];

  elements.forEach(element => {
    if (element.type === 'checkbox' && element.checked) {
      completedIds.push(element.id);
    }
  });

  var res = await putCompletedItems(completedIds)
  if (res.status != 204) {
    showErrorPopup('failed to save changes, try again!')
  }

  await reloadList();
}

const registerSearchField = (searchField) => {
  // Add a keypress event listener to the text field.
  searchField.addEventListener('keydown', event => {
    // get search value
    var searchText = searchField.value.toLowerCase()
    var elements = document.querySelectorAll('span')
    if (resetStyle == null) {
      resetStyle = elements[0].parentNode.parentNode.style.display
    }
    elements.forEach( element => {
      if (!element.innerText.toLowerCase().includes(searchText)) {
        element.parentNode.parentNode.classList.add('hidden');
      }
      else {
        element.parentNode.parentNode.classList.remove('hidden');
      }
    });
  });
}

const reloadList = async () => {
  var content = document.getElementById('content');
  content.innerHTML = listDom.getElementById('list-template').innerHTML;

  var listData = await fetch(fullPath);
  var listJson = await listData.json()

  var listTemplate = listDom.getElementById("list-item-template");
  listJson.forEach(element => {
    if (!element.done) {
      listTemplate.content.querySelector("span").innerText = element.item
      listTemplate.content.querySelector("input").id = element.id
      var node = document.importNode(listTemplate.content.querySelector("div"), true)
      document.getElementById('list-container').appendChild(node)
    }
  });
}

// main entrypoint

document.addEventListener('DOMContentLoaded', async () => {
  var content = document.getElementById('content');

  // get the base path
  if (basePath == null) {
    basePath = new URL(window.location.href).origin;
    fullPath = basePath + endpoint;
  }

  // load the list content
  listDom = await loadAndCacheHtmlFile('list.html');
  await reloadList();

  registerSearchField(document.getElementById('search-box'));
});