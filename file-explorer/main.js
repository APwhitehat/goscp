const {app, BrowserWindow, ipcMain} = require('electron');
const electron = require('electron');
const node_ssh = require('node-ssh');
const fs = require('fs');

function fs_writeFile(fileName, data) {
	return new Promise(function(resolve, reject) {
		fs.writeFile(fileName, data, 'utf-8', function(err) {
			if (err) reject(err);
			else resolve(data);
		});
	});
}
function fs_readFile(fileName, data) {
	return new Promise(function(resolve, reject) {
		fs.readFile(fileName, function(err, data) {
			if (err) reject(err);
			else resolve(data);
		});
	});
}

var ssh = new node_ssh();
// Enable live reload for Electron too
require('electron-reload')(__dirname, {
	// Note that the path to electron may vary according to the main file
	electron: require(`${__dirname}/node_modules/electron`)
});
let mainWindow;

// Quit when all windows are closed.
app.on('window-all-closed', function() {
	if (process.platform != 'darwin')
		app.quit();
});

// This method will be called when Electron has done everything
// initialization and ready for creating browser windows.
app.on('ready', function() {
	// Create the browser window.
	mainWindow = new BrowserWindow({width: 800, height: 600});
	credentialsFile = "./credentials.txt";

	/*mainWindow.webContents.openDevTools();*/

	// and load the index.html of the app.
	mainWindow.loadURL('file://' + __dirname + '/index.html');

	mainWindow.webContents.once('dom-ready', () => {
		fs_readFile(credentialsFile)
			.then((data) => JSON.parse(data))
			.then((data) => {
				console.log(data);
				mainWindow.webContents.send('restoreForm', data);
				console.log("Restored Form Values");
			})
			.catch(console.error);


	})
	// Emitted when the window is closed.
	mainWindow.on('closed', function() {
		// Dereference the window object, usually you would store windows
		// in an array if your app supports multi windows, this is the time
		// when you should delete the corresponding element.
		mainWindow = null;
	});



	ipcMain.on('getDirectory', (evt, data) => {
		console.log(data);
		delete(data.password);
		if (data.saveValue) {
			fs_writeFile(credentialsFile, JSON.stringify(data))
				.then((val) => console.log("Write Successful"))
				.catch(console.error);
		}
	})
	/*  ssh.connect({
	 *    host: 'localhost',
	 *    username: 'msharma',
	 *    privateKey: '/home/msharma/.ssh/msharma'
	 *  })
	 *    .then(() => {
	 *      let homeDir = '/home/msharma/'
	 *      ssh.getFile('/home/msharma/Downloads/index.html',`${homeDir}/Public/foo` )	
	 *        .then((contents)=> {
	 *
	 *          console.log("The File's contents were successfully downloaded");
	 *        })
	 *      ssh.execCommand('ls -la', {
	 *        cwd: `${homeDir}/Public` 
	 *      })
	 *        .then((result) => {
	 *          console.log(result.stdout.split('\n'));
	 *        })
	 *        .catch((err) => console.error(err))
	 *    })
	 *    .catch((err) => console.error(err))*/
});
