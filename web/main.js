var listDom;
var basePath;
var resetStyle;
var fullPath;

const endpoint = '/api/v1/list/';

const showErrorPopup = (errorMessage, statusCode) => {
  // Set the error message.
  document.getElementById('error-message').innerHTML = `${errorMessage} HTTP: ${statusCode}`;
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
  return await fetch(request);
}

// provide a list of items which are updated to done
const putCompletedItems = async (itemList, list) => {
  // Create a new Fetch request object.
  const request = new Request(fullPath, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(itemList),
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

const saveNew = async (listId) => {
  // take the listId from the query path if not supplied?
  if (listId == null) {
    var url = new URL(window.location.href)
    listId = url.searchParams.get("listId")
  }
  var itemText = document.getElementById('new-item-text').value
  var res = await postItem(itemText, listId)
  // CHECK IF IT IS A 201 CREATED
  if (res.status != 201) {
    showErrorPopup(`list with ID '${listId}' could not be retrieved`, res.status)
    return
  }

  var container = document.getElementById('list-container')
  container.children[1].remove()
  var listTemplate = listDom.getElementById("list-item-template");
  listTemplate.content.querySelector("span").innerText = itemText
  listTemplate.content.querySelector("input").id = "test-id"
  var node = document.importNode(listTemplate.content.querySelector("div"), true)
  container.insertBefore(node, container.children[1])
}

const saveChanges = async (listId) => {
  //TODO: remove this duplication, should be more elegant
  if (listId == null) {
    var url = new URL(window.location.href)
    listId = url.searchParams.get("listId")
  }
  const elements = document.querySelectorAll('input[type="checkbox"]');
  const completedItems = [];
  elements.forEach(element => {
    if (element.type === 'checkbox' && element.checked) {
      completedItems.push({
        id: element.id,
        listId: listId,
        done: true
      });
    }
  });

  var res = await putCompletedItems(completedItems, listId)
  if (res.status != 204) {
    showErrorPopup('failed to save changes, try again!', res.status)
  }

  await loadList(listId);
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

const loadList = async (listId) => {
  if (listId == null) {
    var listId = document.getElementById("list-id-input").value;
  }
  history.pushState(null, null, window.location.pathname + "?listId=" + encodeURIComponent(listId));
  var res = await fetch(fullPath + listId);
  var listJson = await res.json()

  if (res.status != 200 || listJson == null) {
    showErrorPopup(`list with ID '${listId}' could not be retrieved`, res.status);
    return
  }

  var content = document.getElementById('content');
  content.innerHTML = listDom.getElementById('list-template').innerHTML;
  var listTemplate = listDom.getElementById("list-item-template");
  listJson.forEach(element => {
    if (!element.done) {
      listTemplate.content.querySelector("span").innerText = element.item
      listTemplate.content.querySelector("input").id = element.id
      var node = document.importNode(listTemplate.content.querySelector("div"), true)
      document.getElementById('list-container').appendChild(node)
    }
  });
  registerSearchField(document.getElementById('search-box'));
}

// main entrypoint

document.addEventListener('DOMContentLoaded', async () => {
  // get the base path
  var url = new URL(window.location.href)
  if (basePath == null) {
    basePath = url.origin;
    fullPath = basePath + endpoint;
  }

  // load the list content
  listDom = await loadAndCacheHtmlFile('list.html');

  // load the list if specified in the query param
  var listId = url.searchParams.get("listId")
  if (listId != null) {
    loadList(listId)
  }
});