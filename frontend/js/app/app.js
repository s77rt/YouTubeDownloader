const nano = 1/1000000000;

const videos = document.querySelector('#videos');
const video_template = document.querySelector('#video_template');
var video_placeholder;

const error_box = document.querySelector('#error_box');

const toggleLightDark__button = document.querySelector('#toggleLightDark__button');

const before_results = document.querySelector('#before_results');

const getvideo_form = document.querySelector('#getvideo_form');
const getvideo_input = document.querySelector('#getvideo_input');

const info = document.querySelector('#info');

var lock = false;

var DB = {}
DB.Videos = {}

var theme = document.querySelector("body").getAttribute("data-theme") || "light";

////////////////////////////////////////////////////////////////

function create_video_placeholder() {
	remove_video_placeholder();

	var video = video_template.cloneNode(true);

	video.setAttribute('id', "video_placeholder");
	video.style.removeProperty("display");

	video.querySelector('[data-role=title]').innerHTML = "<br><br><br>";
	video.querySelector('[data-role=title]').classList.add("skeleton-box");
	video.querySelector('[data-role=duration]').remove();
	video.querySelector('[data-role=author]').parentElement.classList.add("skeleton-box");
	video.querySelector('[data-role=author]').innerHTML = "<br>";
	video.querySelector('[data-role=thumbnail]').parentElement.classList.add("skeleton-box");
	video.querySelector('[data-role=thumbnail]').remove();
	video.querySelector('[data-role=remove]').parentElement.classList.add("skeleton-box");
	video.querySelector('[data-role=remove]').style.opacity = 0;
	video.querySelector('[data-role=remove]').style.cursor = "inherit";
	video.querySelector('[data-role=info]').remove();
	video.querySelector('[data-role=download]').remove();

	var formats = video.querySelector('[data-role=formats]');
	for (var i = 0; i < 4; i++) {
		var format = document.createElement('button');
		format.classList.add("badge");
		format.classList.add("btn");
		format.classList.add("p-2");
		format.classList.add("skeleton-box");
		format.innerHTML = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;";
		formats.append(format);
	}

	before_results_hide();

	video_placeholder = video;
	videos.prepend(video_placeholder);

	new ScrollBooster({
		viewport: formats.parentElement,
		content: formats,
		scrollMode: 'transform',
		direction: 'horizontal',
		emulateScroll: true,
	});
}

function remove_video_placeholder() {
	if (video_placeholder != null)
		video_placeholder.remove();
	before_results_auto();
}

////////////////////////////////////////////////////////////////

function add_video(video_obj) {
	var video = video_template.cloneNode(true);
	let video_uuid = 'v'+generateUUID();

	DB.Videos[video_uuid] = {"status": "idle", "info": video_obj};

	video.setAttribute('id', video_uuid);
	video.style.removeProperty("display");

	video.ontransitionend = (e) => {
		if (e.srcElement.classList.contains('toberemoved')) {
			e.srcElement.remove();
			delete(DB.Videos[video_uuid]);
			onDBdelete();
		}
	};

	video.querySelector('[data-role=title]').innerText = video_obj.Title;
	video.querySelector('[data-role=duration]').innerText = format_time(video_obj.Duration*nano);
	video.querySelector('[data-role=author]').innerText = video_obj.Author;
	video.querySelector('[data-role=thumbnail]').src = get_best_thumbnail(video_obj.Thumbnails).URL;

	var formats = video.querySelector('[data-role=formats]');
	var added_formats = [];
	for (var i = 0; i < video_obj.Formats.length; i++) {
		let format_obj = video_obj.Formats[i];

		if (format_obj.mimeType.startsWith("audio/"))
			format_obj.qualityLabel = "Audio"

		if (format_obj.qualityLabel == "")
			continue;

		if (added_formats.find(f => f.qualityLabel === format_obj.qualityLabel))
			continue;

		if (format_obj.height == 0) {
			format_btn_style = "btn-outline-secondary";
		} else if (format_obj.height > 1080) {
			format_btn_style = "btn-success";
		} else if (format_obj.height >= 720) {
			format_btn_style = "btn-danger";
		} else if (format_obj.height >= 480) {
			format_btn_style = "btn-primary";
		} else {
			format_btn_style = "btn-secondary";
		}

		var format = document.createElement('button');
		format.classList.add("badge");
		format.classList.add("btn");
		format.classList.add("p-2");
		format.classList.add(format_btn_style);
		format.innerText = format_obj.qualityLabel;
		format.addEventListener('click', function() {
			w_DownloadVideo(video_uuid, video_obj, format_obj);
		});

		format.setAttribute("title", humanFileSize(format_obj.contentLength))

		formats.append(format);

		added_formats.push(format_obj);
	}

	video.querySelector('[data-role=download]').addEventListener('click', function() {
		video.querySelectorAll('[data-role=formats]').forEach(function(f) {
			let best_format = f.firstElementChild;
			if (best_format != null)
				best_format.click();
		});
	});
	video.querySelector('[data-role=info]').addEventListener('click', function() {
		let f = new JSONFormatter(DB.Videos[video_uuid].info, 1, {
			theme: theme,
		});
		info.querySelector('[data-role=info_json').replaceChildren(f.render());
	});
	video.querySelector('[data-role=remove]').addEventListener('click', function() {
		remove_video(video_uuid);
	});

	video.querySelector('[data-role=cancel]').addEventListener('click', function() {
		w_CancelDownload(video_uuid);
	});
	video.querySelector('[data-role=play]').addEventListener('click', function() {
		w_PlayVideo(video_uuid);
	});
	video.querySelector('[data-role=openfolder]').addEventListener('click', function() {
		w_OpenFolder(video_uuid);
	});
	video.querySelector('[data-role=delete]').addEventListener('click', function() {
		w_DeleteVideo(video_uuid);
	});

	videos.prepend(video);

	new ScrollBooster({
		viewport: formats.parentElement,
		content: formats,
		scrollMode: 'transform',
		direction: 'horizontal',
		emulateScroll: true,
	});

	update_video_status(video_uuid, "idle");

	w_GetExistingVideo(video_uuid, video_obj, added_formats);

	onDBinsert();
}

function remove_video(video_uuid) {
	if (DB.Videos[video_uuid].status != "idle" && DB.Videos[video_uuid].status != "downloaded")
		return

	var video = document.querySelector('#'+video_uuid);

	video.classList.add("toberemoved");
	video.classList.remove("show");

	// also see video.ontransitionend
}

function update_video_status(video_uuid, status) {
	DB.Videos[video_uuid].status = status;

	var video = document.querySelector('#'+video_uuid);

	switch (status) {
		case "fetching":
			video_clearAlert(video_uuid);

			// no ui element change is required here. (fetching is kind of a silent task)

			break;
		case "downloading":
			video_clearAlert(video_uuid);

			video.querySelector('[data-role=progressbar]').style.width = "0%";
			video.querySelector('[data-role=progressbar]').innerText = "";
			video.querySelector('[data-role=progressbar]').classList.remove("progress-bar-striped");
			video.querySelector('[data-role=progressbar]').classList.remove("progress-bar-animated");
			video.querySelector('[data-role=statusbar]').innerHTML = '<i class="spinner-border spinner-border-sm"></i>';
			video.querySelector('[data-role=speed]').innerText = "";

			video.querySelector('[data-role=download_options]').style.setProperty("display", "none", "important");
			video.querySelector('[data-role=download_status]').style.removeProperty("display");
			video.querySelector('[data-role=actions]').style.setProperty("display", "none", "important");

			video.querySelector('[data-role=remove]').classList.add("disabled");
			video.querySelector('[data-role=download]').classList.add("disabled");

			break;
		case "downloaded":
			video.querySelector('[data-role=progressbar]').style.width = "100%";
			video.querySelector('[data-role=progressbar]').innerText = "Downloaded 100%";
			video.querySelector('[data-role=progressbar]').classList.remove("progress-bar-striped");
			video.querySelector('[data-role=progressbar]').classList.remove("progress-bar-animated");
			video.querySelector('[data-role=statusbar]').innerText = "Downloaded 100%";
			video.querySelector('[data-role=speed]').innerText = "";

			video.querySelector('[data-role=download_options]').style.setProperty("display", "none", "important");
			video.querySelector('[data-role=download_status]').style.setProperty("display", "none", "important");
			video.querySelector('[data-role=actions]').style.removeProperty("display");

			video.querySelector('[data-role=remove]').classList.remove("disabled");
			video.querySelector('[data-role=download]').classList.add("disabled");

			break;
		case "idle":
			video.querySelector('[data-role=progressbar]').style.width = "0%";
			video.querySelector('[data-role=progressbar]').innerText = "";
			video.querySelector('[data-role=progressbar]').classList.remove("progress-bar-striped");
			video.querySelector('[data-role=progressbar]').classList.remove("progress-bar-animated");
			video.querySelector('[data-role=statusbar]').innerText = "";
			video.querySelector('[data-role=speed]').innerText = "";

			video.querySelector('[data-role=download_options]').style.removeProperty("display");
			video.querySelector('[data-role=download_status]').style.setProperty("display", "none", "important");
			video.querySelector('[data-role=actions]').style.setProperty("display", "none", "important");

			video.querySelector('[data-role=remove]').classList.remove("disabled");
			video.querySelector('[data-role=download]').classList.remove("disabled");

			break;
	}
}

function update_video_progress(video_uuid, progress_1_obj, progress_2_obj) {
	if (DB.Videos[video_uuid].status != "downloading")
		return

	var video = document.querySelector('#'+video_uuid);

	let percentage, percentage_text, percentage_animated, status_text, speed;

	var coefficients = {
		"DownloadingVideo": 0.56,
		"DownloadingAudio": 0.37,
		"Merging": 0.07,
	}

	function local_percentage(offset, transferred, size, coefficient) {
		let percentage;
		percentage = offset;
		if (size != 0)
			percentage += (transferred / size)*coefficient;
		percentage *= 100;
		return percentage;
	}

	if (progress_1_obj.Size == 0 || progress_1_obj.Transferred != progress_1_obj.Size) {
		// Step1: Downloading Part1
		percentage = local_percentage(
			0,
			progress_1_obj.Transferred,
			progress_1_obj.Size,
			coefficients.DownloadingVideo
		).toFixed(2);
		percentage_text = "1/3 Downloading Part1" + " (" + percentage + "%)";
		percentage_animated = false;
		status_text = "1/3 Downloading Part1" + " (" + percentage + "%)";
		speed = humanFileSize((progress_1_obj.Speed)/(progress_1_obj.TimeUnit*nano))+"/s";
	} else if (progress_2_obj.Size == 0 || progress_2_obj.Transferred != progress_2_obj.Size) {
		// Step2: Downloading Part2
		percentage = local_percentage(
			coefficients.DownloadingVideo,
			progress_2_obj.Transferred,
			progress_2_obj.Size,
			coefficients.DownloadingAudio
		).toFixed(2);
		percentage_text = "2/3 Downloading Part2" + " (" + percentage + "%)";
		percentage_animated = false;
		status_text = "2/3 Downloading Part2" + " (" + percentage + "%)";
		speed = humanFileSize((progress_2_obj.Speed)/(progress_2_obj.TimeUnit*nano))+"/s";
	} else if (progress_1_obj.Transferred == progress_1_obj.Size && progress_2_obj.Transferred == progress_2_obj.Size) {
		// Step3: Merging
		percentage = local_percentage(
			coefficients.DownloadingVideo + coefficients.DownloadingAudio,
			0,
			0,
			coefficients.Merging
		).toFixed(2);
		percentage_text = "3/3 Merging...";
		percentage_animated = true;
		status_text = "3/3 Merging...";
		speed = "N/A";
	} else {
		// Unknown progress
		percentage = 100;
		percentage_text = "Downloading... (Unknown percentage)";
		percentage_animated = true;
		status_text = "Downloading... (Unknown percentage)";
		speed = "N/A";
	}

	video.querySelector('[data-role=progressbar]').style.width = percentage+"%";
	video.querySelector('[data-role=progressbar]').innerText = percentage_text;
	if (percentage_animated) {
		video.querySelector('[data-role=progressbar]').classList.add("progress-bar-striped");
		video.querySelector('[data-role=progressbar]').classList.add("progress-bar-animated");
	} else {
		video.querySelector('[data-role=progressbar]').classList.remove("progress-bar-striped");
		video.querySelector('[data-role=progressbar]').classList.remove("progress-bar-animated");
	}
	video.querySelector('[data-role=statusbar]').innerText = status_text;
	video.querySelector('[data-role=speed]').innerText = speed;
}

function update_video_success(video_uuid, video_filename) {
	if (DB.Videos[video_uuid].status != "downloading" && DB.Videos[video_uuid].status != "fetching")
		return

	update_video_status(video_uuid, "downloaded");

	DB.Videos[video_uuid].filename = video_filename;
}

////////////////////////////////////////////////////////////////

function newAlert(origin, msg) {
	var alert = document.createElement("div");
	alert.classList.value = "alert alert-danger alert-dismissible fade show";

	var alert_msg = document.createElement("div");
	alert_msg.innerHTML = '<strong><i class="bi-exclamation-circle-fill"></i> '+origin+'</strong> '+msg;

	var close_btn = document.createElement("button");
	close_btn.classList.value = "btn-close";
	close_btn.setAttribute("type", "button");
	close_btn.setAttribute("data-bs-dismiss", "alert");
	close_btn.setAttribute("aria-label", "Close");

	alert.appendChild(alert_msg);
	alert.appendChild(close_btn);

	error_box.appendChild(alert);
}

function video_newAlert(video_uuid, origin, msg) {
	var video = document.querySelector('#'+video_uuid);

	var error_box = video.querySelector('[data-role=error_box]');

	var alert = document.createElement("div");
	alert.classList.value = "alert alert-danger alert-dismissible fade show";

	var alert_msg = document.createElement("div");
	alert_msg.innerHTML = '<strong><i class="bi-exclamation-circle-fill"></i> '+origin+'</strong> '+msg;

	var close_btn = document.createElement("button");
	close_btn.classList.value = "btn-close";
	close_btn.setAttribute("type", "button");
	close_btn.setAttribute("data-bs-dismiss", "alert");
	close_btn.setAttribute("aria-label", "Close");

	alert.appendChild(alert_msg);
	alert.appendChild(close_btn);

	error_box.replaceChildren(alert);
}

function video_clearAlert(video_uuid) {
	var video = document.querySelector('#'+video_uuid);

	var error_box = video.querySelector('[data-role=error_box]');

	error_box.replaceChildren();
}

////////////////////////////////////////////////////////////////

function w_CheckFFmpeg() {
	CheckFFmpeg();
}

function w_CheckFFmpeg__onerror(error) {
	newAlert("FFmpeg", error);
}

////////////////////////////////////////////////////////////////

function w_GetVideo(url) {
	GetVideo(url);
}

function w_GetVideo__onerror(error) {
	newAlert("GetVideo", error);
}

function w_GetVideos(url) {
	GetVideos(url);
}

function w_GetVideos__onerror(error) {
	newAlert("GetVideos", error);
}

function w_GetExistingVideo(video_uuid, video_obj, format_objs) {
	if (DB.Videos[video_uuid].status != "idle")
		return

	update_video_status(video_uuid, "fetching");

	GetExistingVideo(video_uuid, video_obj, format_objs);
}
function w_GetExistingVideo__onerror(video_uuid, error) {
	if (error != "not available")
		video_newAlert(video_uuid, "GetExistingVideo", error);
	update_video_status(video_uuid, "idle");
}

function w_DownloadVideo(video_uuid, video_obj, format_obj) {
	if (DB.Videos[video_uuid].status != "idle")
		return

	update_video_status(video_uuid, "downloading");
	
	DownloadVideo(video_uuid, video_obj, format_obj);
}
function w_DownloadVideo__onerror(video_uuid, error) {
	update_video_status(video_uuid, "idle");
	video_newAlert(video_uuid, "DownloadVideo", error);
}

function w_CancelDownload(video_uuid) {
	if (DB.Videos[video_uuid].status != "downloading")
		return

	update_video_status(video_uuid, "idle");

	CancelDownload(video_uuid);
}
function w_CancelDownload__onerror(video_uuid, error) {
	update_video_status(video_uuid, "idle");
	video_newAlert(video_uuid, "CancelDownload", error);
}

function w_PlayVideo(video_uuid) {
	if (DB.Videos[video_uuid].status != "downloaded")
		return

	PlayVideo(video_uuid, DB.Videos[video_uuid].filename);
}
function w_PlayVideo__onerror(video_uuid, error) {
	video_newAlert(video_uuid, "PlayVideo", error);
}

function w_OpenFolder(video_uuid) {
	if (DB.Videos[video_uuid].status != "downloaded")
		return

	OpenFolder(video_uuid, DB.Videos[video_uuid].filename);
}
function w_OpenFolder__onerror(video_uuid, error) {
	video_newAlert(video_uuid, "OpenFolder", error);
}

function w_DeleteVideo(video_uuid) {
	if (DB.Videos[video_uuid].status != "downloaded")
		return

	update_video_status(video_uuid, "idle");

	DeleteVideo(video_uuid, DB.Videos[video_uuid].filename);
}
function w_DeleteVideo__onerror(video_uuid, error) {
	video_newAlert(video_uuid, "DeleteVideo", error);
}

////////////////////////////////////////////////////////////////

function before_results_hide() {
	before_results.style.display = "none";
}

function before_results_show() {
	before_results.style.display = "block";
}

function before_results_auto() {
	if (Object.keys(DB.Videos).length > 0) {
		before_results.style.display = "none";
	} else {
		before_results.style.display = "block";
	}
}


function onDBdelete() {
	before_results_auto();
}

function onDBinsert() {
	before_results_auto();
}

////////////////////////////////////////////////////////////////

function get_best_thumbnail(thumbnails) {
	if (thumbnails.length == 0) {
		return {
			"URL": "",
			"Width": 0,
			"Height": 0,
		}
	}

	thumbnail = thumbnails[0];
	for (var i = 1; i < thumbnails.length; i++) {
		thumbnail_tmp = thumbnails[i];
		if (thumbnail_tmp.Height > thumbnail.Height) {
			thumbnail = thumbnail_tmp;
		}
	}

	return thumbnail
}

function format_time(time) {   
	// Hours, minutes and seconds
	var hrs = ~~(time / 3600);
	var mins = ~~((time % 3600) / 60);
	var secs = ~~time % 60;

	// Output like "1:01" or "4:03:59" or "123:03:59"
	var ret = "";
	if (hrs > 0) {
		ret += "" + hrs + ":" + (mins < 10 ? "0" : "");
	}
	ret += "" + mins + ":" + (secs < 10 ? "0" : "");
	ret += "" + secs;
	return ret;
}

function humanFileSize(bytes, si=true, dp=1) {
  const thresh = si ? 1000 : 1024;

  if (Math.abs(bytes) < thresh) {
	return bytes + ' B';
  }

  const units = si 
	? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'] 
	: ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
  let u = -1;
  const r = 10**dp;

  do {
	bytes /= thresh;
	++u;
  } while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);


  return bytes.toFixed(dp) + ' ' + units[u];
}

////////////////////////////////////////////////////////////////

function toggleLightDark() {
	if (theme == "light") {
		theme = "dark";
		document.querySelector("body").setAttribute("data-theme", "dark") 
		toggleLightDark__button.innerHTML = '<i class="bi-moon"></i>';
	} else {
		theme = "light";
		document.querySelector("body").setAttribute("data-theme", "light") 
		toggleLightDark__button.innerHTML = '<i class="bi-sun"></i>';
	}

	let json_formatter = document.querySelector('.json-formatter-open');
	if (json_formatter != null) {
		if (theme == "dark") {
			json_formatter.classList.remove('json-formatter-light');
			json_formatter.classList.add('json-formatter-dark');
		} else {
			json_formatter.classList.remove('json-formatter-dark');
			json_formatter.classList.add('json-formatter-light');
		}
	}
}
toggleLightDark__button.addEventListener('click', toggleLightDark);

function setVersion(version) {
	document.querySelector("#version").innerText = version;
}

////////////////////////////////////////////////////////////////

function generateUUID() { // Public Domain/MIT
	var d = new Date().getTime();//Timestamp
	var d2 = (performance && performance.now && (performance.now()*1000)) || 0;//Time in microseconds since page-load or 0 if unsupported
	return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
		var r = Math.random() * 16;//random number between 0 and 16
		if(d > 0){//Use timestamp until depleted
			r = (d + r)%16 | 0;
			d = Math.floor(d/16);
		} else {//Use microseconds since page-load if supported
			r = (d2 + r)%16 | 0;
			d2 = Math.floor(d2/16);
		}
		return (c === 'x' ? r : (r & 0x3 | 0x8)).toString(16);
	});
}

var STRIP_COMMENTS = /((\/\/.*$)|(\/\*[\s\S]*?\*\/))/mg;
var ARGUMENT_NAMES = /([^\s,]+)/g;
function getParamNames(func) {
	var fnStr = func.toString().replace(STRIP_COMMENTS, '');
	var result = fnStr.slice(fnStr.indexOf('(')+1, fnStr.indexOf(')')).match(ARGUMENT_NAMES);
	if(result === null)
		result = [];
	return result;
}

////////////////////////////////////////////////////////////////

function h_lock() {
	lock = true;
	getvideo_input.setAttribute("disabled", "disabled");
	create_video_placeholder();
}

function h_unlock() {
	lock = false;
	getvideo_input.removeAttribute("disabled");
	remove_video_placeholder();
}

////////////////////////////////////////////////////////////////

function downloadAll() {
	document.querySelectorAll('[data-role=download]').forEach(function(f) {
		f.click();
	});
}

function cancelAll() {
	document.querySelectorAll('[data-role=cancel]').forEach(function(f) {
		f.click();
	});
}

function removeAll() {
	document.querySelectorAll('[data-role=remove]').forEach(function(f) {
		f.click();
	});
}

function clear() { removeAll(); }

////////////////////////////////////////////////////////////////

document.onkeydown = function (event) {
	if (event.ctrlKey) { // Ctrl is pressed
		// Allown only A, X, C, V
		if (event.keyCode != 65 && event.keyCode != 88 && event.keyCode != 67 && event.keyCode != 86)
			return false;
	}
	if (event.keyCode >= 112 && event.keyCode <= 123) // F1 to F12
		return false;
}
window.addEventListener("contextmenu", function(e) { e.preventDefault(); })

getvideo_form.addEventListener("submit", function(e) {
	e.preventDefault();
	if (lock)
		return
	h_lock();
	GetVideos(getvideo_input.value);
});

document.querySelector("#control_downloadall").addEventListener("click", function() {
	downloadAll();
});
document.querySelector("#control_cancelall").addEventListener("click", function() {
	cancelAll();
});
document.querySelector("#control_clear").addEventListener("click", function() {
	clear();
});

// Anchors default behaviour patch
document.querySelector('body').addEventListener('click', function(e) {
	if (e.target.tagName == "A") {
		e.preventDefault();
		if (e.target.href != "")
			window.open(e.target.href, '_blank').focus();
	}
}, true);
