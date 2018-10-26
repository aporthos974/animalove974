function notifySuccess(title, message) {
	notify('success', title, message);
}

function notifyInfo(title, message) {
	notify('info', title, message, 4000);
}

function notifyWarning(title, message) {
	notify('warning', title, message);
}

function notifyError(title, message) {
	notify('danger', title, message);
}

function notify(type, title, message, delay) {
	$.notify({
		// options
		icon: 'glyphicon glyphicon-warning-sign',
		title: '<strong>'+title+'</strong><br/>',
		message: message
	},{
		// settings
		element: 'body',
		position: null,
		type: type,
		allow_dismiss: true,
		newest_on_top: false,
		placement: {
			from: "top",
			align: "right"
		},
		offset: 20,
		spacing: 10,
		z_index: 5031,
		delay: delay ? delay : 6000,
		timer: 1000,
		url_target: '_blank',
		animate: {
			enter: 'animated fadeInDown',
			exit: 'animated fadeOutUp'
		},
		icon_type: 'class',
		template: '<div id="notification-popin" data-notify="container" class="col-xs-11 col-sm-3 alert alert-{0}" role="alert">' +
			'<button type="button" aria-hidden="true" class="close" data-notify="dismiss">×</button>' +
			'<span data-notify="icon"></span> ' +
			'<span data-notify="title">{1}</span> ' +
			'<span data-notify="message">{2}</span>' +
		'</div>' 
	});
}