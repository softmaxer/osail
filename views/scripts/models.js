const models = [];
let start_port = 11435;

function addModel() {
  modelName = document.getElementById("model-name").value;
  ;
  let element = {"name": modelName, "host": "http://localhost", port: start_port};
  models.push(element);
  start_port += 1;
  displayModels(element);
}

function displayModels(element) {
  modelsList = document.getElementById("models");
  let llmNode = document.createElement('li');
  llmNode.innerHTML = element["name"];
  modelsList.appendChild(llmNode);
  
}


function startExperiment() {
  
}
