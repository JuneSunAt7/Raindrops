document.getElementById('saveBtn').addEventListener('click', function () {
    var text = document.getElementById('textarea').value;
    
    var blob = new Blob([text], {type: 'text/plain'});
    var url = window.URL.createObjectURL(blob);
    
    var a = document.createElement('a');
    a.href = url;
    a.download = 'textarea_content.txt';
    document.body.appendChild(a);
    a.click();
    
    window.URL.revokeObjectURL(url);
  });

  document.getElementById('openBtn').addEventListener('click', function () {
    var input = document.createElement('input');
    input.type = 'file';

    input.onchange = e => { 
        var file = e.target.files[0]; 
        var reader = new FileReader();

        reader.onload = function() {
            var text = reader.result;
            document.getElementById('textarea').value = text;
        };

        reader.readAsText(file);
    };

    input.click();
});
