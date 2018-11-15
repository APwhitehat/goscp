global.$ = $;
const electron = require('electron');
const remote  = electron.remote;
const {Menu, BrowserWindow, MenuItem, shell} = remote;
const ipc = electron.ipcRenderer;


function RadioOnClickHandler() {
	let form = document.getElementById('form');
	[form['password'].disabled, form['privKeyPath'].disabled] = [form['privKeyPath'].disabled, form['password'].disabled] ;
}
(function() {
  'use strict';
  window.addEventListener('load', function() {
    // Fetch all the forms we want to apply custom Bootstrap validation styles to
		var forms = document.getElementsByClassName('needs-validation');
		let rbutton1 = document.getElementById('rbutton1');
		let rbutton2 = document.getElementById('rbutton2');
		rbutton1.onclick = RadioOnClickHandler;
		rbutton2.onclick = RadioOnClickHandler;

		ipc.on('restoreForm', (evt, data) => {
			console.log(data);
			let form = document.getElementById('form');
			if (data.hostname != undefined)
				form['hostname'].value = data.hostname;
			form['privKeyPath'].value = data.privKeyPath;
			if (data.username != undefined)
				form['username'].value = data.username;
			if (data.src != undefined)
				form['src'].value = data.src;
			if (data.dest != undefined)
				form['dest'].value = data.dest;

			form['checkbox'].checked = data.saveValue;
		})
    // Loop over them and prevent submission
    var validation = Array.prototype.filter.call(forms, function(form) {
      form.addEventListener('submit', function(event) {
        if (form.checkValidity() === false) {
          event.preventDefault();
          event.stopPropagation();
        }
        form.classList.add('was-validated');
      }, false);
    });
  }, false);
})();
// append default actions to menu for OSX
function formSubmit(){
	setTimeout(()=> {console.log("HIIII")}, 2000);
	let form = document.getElementById('form');
	console.log(form);
	let formData = {
		hostname: form['hostname'].value,
		username: form['username'].value,
		privKeyPath: form['privKeyPath'].value,
		password: form['password'].value,
		src: form['src'].value,
		dest: form['dest'].value,
		saveValue: form['checkbox'].checked
	}
	console.log(formData);
	/*alert('Done!');*/
	if (formData.hostname == "" || formData.username == "" || formData.src == "" || formData.dest == ""){
		console.log("Form not filled!");
		return;
	}
	ipc.send('getDirectory', formData);
	let endPoint = 'http://localhost:8080/scp';

	fetch(endPoint, {
        method: "POST", // *GET, POST, PUT, DELETE, etc.
        mode: "cors", // no-cors, cors, *same-origin
        cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
        credentials: "same-origin", // include, same-origin, *omit
        headers: {
            "Content-Type": "application/json; charset=utf-8",
        },
        body: JSON.stringify(formData), // body data type must match "Content-Type" header
	})
		.then((val) => {
			alert("Files are being transferred!");
			console.log(val);
		})
		.catch((err) => {
			console.error;
			alert("Failed!");
		});
}
var initMenu = function () {
	try {
		var nativeMenuBar = new Menu();
		if (process.platform == "darwin") {
			nativeMenuBar.createMacBuiltin && nativeMenuBar.createMacBuiltin("FileExplorer");
		}
	} catch (error) {
		console.error(error);
		setTimeout(function () { throw error }, 1);
	}
};

var aboutWindow = null;
var credentialsWindow = null;
var App = {
	// show "about" window
	about: function () {
		var params = {toolbar: false, resizable: false, show: true, height: 150, width: 400};
		aboutWindow = new BrowserWindow(params);
		aboutWindow.loadURL('file://' + __dirname + '/about.html');
	},
	// change folder for sidebar links
	cd: function (anchor) {
		anchor = $(anchor);

		$('#sidebar li').removeClass('active');
		$('#sidebar i').removeClass('icon-white');

		anchor.closest('li').addClass('active');
		anchor.find('i').addClass('icon-white');

		this.setPath(anchor.attr('nw-path'));
	},

	// set path for file explorer
	setPath: function (path) {
		if (path.indexOf('~') == 0) {
			path = path.replace('~', process.env['HOME']);
		}
		this.folder.open(path);
		this.addressbar.set(path);
	}
};

$(document).ready(function() {
	initMenu();
});
