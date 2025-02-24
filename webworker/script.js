
const drake = dragula({
    containers: [
        document.getElementById("to-do"),
        document.getElementById("doing"),
        document.getElementById("done"),
        document.getElementById("trash")
    ],
    removeOnSpill: false 
});

function addTask() {

    var inputTask = document.getElementById("taskText").value;

    if (inputTask != '') {
        document.getElementById("to-do").innerHTML += "<li class='task'><p>" + inputTask + "</p></li>";

        document.getElementById("taskText").value = "";
    }
}

function emptyTrash() {
    document.getElementById("trash").innerHTML = "";
    saveTasks();
}

document.addEventListener('DOMContentLoaded', (event) => {
    loadTasks(); 

    drake.on('drop', saveTasks); 
});

function loadTasks() {
    const columns = ['to-do', 'doing', 'done', 'trash'];
    
    columns.forEach(column => {
        const taskList = document.getElementById(column);
        taskList.innerHTML = ''; 
        const tasks = JSON.parse(localStorage.getItem(column)) || [];
        tasks.forEach(taskText => addTaskToColumn(taskText, column));
    });
}

function saveTasks() {
    const columns = ['to-do', 'doing', 'done', 'trash'];
    
    columns.forEach(column => {
        const tasks = Array.from(document.getElementById(column).children)
                           .map(task => task.querySelector('p').textContent);
        localStorage.setItem(column, JSON.stringify(tasks));
    });
}

function addTaskToColumn(taskText, columnId) {
    const taskList = document.getElementById(columnId);
    const li = document.createElement('li');
    li.className = 'task';
    li.innerHTML = `<p>${taskText}</p>`;
    taskList.appendChild(li);
}

function emptyTrash() {
    const trash = document.getElementById('trash');
    trash.innerHTML = ''; 
    saveTasks(); 
}