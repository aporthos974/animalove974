"use strict"

var app = angular.module('reunion', ['mapRaphael', 'ui.select', 'ngSanitize', 'ui.bootstrap']);

app.controller('ReunionCtrl', ['$rootScope', '$scope', '$http', '$location', '$q', 'WSService', 'utils', function($rootScope, $scope, $http, $location, $q, WSService, utils) {
	
	$scope.cities = [ { id: 'stdenis', label: 'Saint Denis', region: 'Nord', img: '/static/images/map.png' }, { id: 'standre', label: 'Saint André', region: 'Est', img: '/static/images/map.png' }, { id: 'stpaul', label: 'Saint Paul', region: 'Ouest', img: '/static/images/map.png' },{ id: 'lesavirons', label: 'Les Avirons', region: 'Sud', img: '/static/images/map.png' },
					{ id: 'stemarie', label: 'Sainte Marie', region: 'Nord', img: '/static/images/map.png' }, { id: 'lapossession', label: 'La Possession', region: 'Ouest', img: '/static/images/map.png' }, { id: 'stesuzanne', label: 'Sainte Suzanne', region: 'Nord', img: '/static/images/map.png' },
					{ id: 'leport', label: 'Le Port', region: 'Ouest', img: '/static/images/map.png' }, { id: 'stleu', label: 'Saint Leu', region: 'Ouest', img: '/static/images/map.png' }, { id: 'stlouis', label: 'Saint Louis', region: 'Sud', img: '/static/images/map.png' }, { id: 'stpierre', label: 'Saint Pierre', region: 'Sud', img: '/static/images/map.png' }, { id: 'letampon', label: 'Le Tampon', region: 'Sud', img: '/static/images/map.png' },
					{ id: 'entredeux', label: 'Entre Deux', region: 'Sud', img: '/static/images/map.png' }, { id: 'cilaos', label: 'Cilaos', region: 'Sud', img: '/static/images/map.png' }, { id: 'salazie', label: 'Salazie', region: 'Est', img: '/static/images/map.png' }, { id: 'laplainedespalmistes', label: 'La Plaine des Palmistes', region: 'Est', img: '/static/images/map.png' },
					{ id: 'sterose', label: 'Sainte Rose', region: 'Est', img: '/static/images/map.png' }, { id: 'braspanon', label: 'Bras Panon', region: 'Est', img: '/static/images/map.png' }, { id: 'stbenoit', label: 'Saint Benoit', region: 'Est', img: '/static/images/map.png' }, { id: 'stjoseph', label: 'Saint Joseph', region: 'Est', img: '/static/images/map.png' }, { id: 'stphilippe', label: 'Saint Philippe', region: 'Est', img: '/static/images/map.png' },
					{ id: 'troisbassins', label: 'Trois Bassins', region: 'Ouest', img: '/static/images/map.png' }, { id: 'petiteile', label: 'Petite ile', region: 'Sud', img: '/static/images/map.png' }, { id: 'etangsale', label: 'Étang Salé', region: 'Sud', img: '/static/images/map.png' } ];

	$rootScope.numberByPage = 5;
		
	$scope.oneResultTemplates = { 'perdu': '/static/html/loss_one_result.html', 'errant': '/static/html/seen_one_result.html', 'adopter': '/static/html/adopt_one_result.html', 'found': '/static/html/found_one_result.html' };
	
	
	$scope.getOneResultTemplate = function(announcement) {
		var path = announcement.State == 'found' ? announcement.State : announcement.Type; 
		return $scope.oneResultTemplates[path];
	};
	
	$scope.isTab = function(path) {
		if (path == 'accueil') {
			return new RegExp('^http://.+/.+$').exec(window.location.toString())== null;
		}
		return new RegExp(path + '$').exec(window.location.toString().split('?')[0]) != null;
	};
	
	$scope.hoverCity = function(city) {
		$('path#' + city).attr('fill', '#F9C866');
	};
	
	$scope.outCity = function(city) {
		$('path#' + city).attr('fill', '#F4EAD6');
	};
	
	$rootScope.$on('seen-created', function(event, announcement) {
		$rootScope.$apply(function() {
			$rootScope.createdSeenAnnouncement = announcement;
		});
		showNotification();
	});
	
	$scope.hideNotification = function() {
		$('.notification').addClass('notification-hidden');
	};
	
	$scope.showContact = function() {
		$('#contact-uncompressed').toggleClass('show');
	};
	
	$scope.hideContact = function() {
		$('#contact-uncompressed').toggleClass('show');
	};
	
	$scope.sendMessage = function() {
		
		$http.post('/ws/contact/message', { sender: $scope.contactSender, message: $scope.contactMessage }).success(function(data, status, headers) {
			$scope.contactSender = undefined;		
			$scope.contactMessage = undefined;		
			notifySuccess('Contact', 'Message envoyé');
		}).error(function(data, status, headers) {
			notifyError('Contact', 'Erreur lors l\'envoi du message.<br/> Veuillez réessayer, si cela persiste merci de contacter l\'administrateur');
		});
	};
	
	function showNotification() {
		$('.notification').removeClass('notification-hidden');
	}	
	
	function searchCityLabel(cities, id) {
		for (var index in cities) {
			if (cities[index].id == id) {
				return cities[index];
			}
		}
		return undefined;
	}
	
	$scope.$on('map-mouseover', function (event, city) {
		$('[name="'+city.id+'"] span').addClass('selected-cities');
		$scope.$apply();
	});
	
	$rootScope.$on('location-saved', function() {
		$('#location').modal('hide');	
	});	
	
	$('#authentication').on('shown.bs.modal', function () {
    	$('#email').focus()
  	});
	$('#contact').on('shown.bs.modal', function () {
    	$('#sender').focus()
  	});
	
	$scope.displayLoading = function (display) {
		$scope.displayProgressbar = display;
	};
	
	
	$scope.selectAnnouncement = function(announcement) {
		$scope.selectedAnnouncement = announcement;
	};
		
	$scope.searchById = function(id) {
		$scope.displayLoading(true);
		$http.get('/ws/animaux?id=' + id + '&offset=0&limit=1').success(function(data, status, headers) {
			$scope.announcements = data.announcements;
			$scope.totalItems = 1;
			$scope.displayLoading(false);
		}).error(function(data, status, headers) {
			notifyError('Recherche', 'Erreur lors de l\'exécution de la recherche.<br/> Veuillez réessayer, si cela persiste merci de contacter l\'administrateur');
			$scope.displayLoading(false);
		});
	};	
	
	$scope.utils = utils;
	$scope.contactSubmitted = false;
	
}]
);


app.factory('WSService', ['$rootScope', function($rootScope) {
	var Service = {};
	var websocket = new WebSocket("ws://" + window.location.host + "/socket");
	
	websocket.onopen = function(){
    };
	
	websocket.onmessage = function(message) {
		var notification = JSON.parse(message.data);
		$rootScope.$emit(notification.action, notification.announcement);
    };
	
	window.onbeforeunload = function() {
	    websocket.close();
	};
	
	return Service;
}]
);

app.service('announcementService', [ '$rootScope', '$q', '$http', function($rootScope, $q, $http) {
	return {
		getAnnouncements: function(offset, type, searchQueryParam) {
			$rootScope.$emit('loading', true);
			var deferred = $q.defer();
			var type = type ? '/' + type : '';
			var searchQueryParam = searchQueryParam ? '&' + searchQueryParam : '';
			$http.get('/ws/animaux' + type + '?offset=' + offset + '&limit=' + $rootScope.numberByPage + searchQueryParam).success(function(data, status, headers) {
				deferred.resolve({ 'announcements': data.announcements, 'totalItems': data.total });
				$rootScope.$emit('loading', false);
			}).error(function(data, status, headers) {
				notifyError('Liste des animaux', 'Erreur lors de la récupération de la liste des animaux.<br/> Veuillez réessayer, si cela persiste merci de contacter l\'administrateur');
				deferred.reject({ 'announcements': [], 'totalItems': 0 });
				$rootScope.$emit('loading', false);
			});
			
			return deferred.promise;
		},
		addLocation: function(announcement, newLocation) {
			if (announcement.Locations == undefined) {
				announcement.Locations = [];
			}
			announcement.Locations.push(newLocation);
			notifyInfo('Localisation', 'Prise en compte de la demande d\'ajout...');
			$http.put('/ws/animaux/perdu/locations', announcement).success(function(data, status, headers) {
	            notifySuccess('Localisation', 'Ajout d\'une nouvelle localisation effectué avec succès');
				$rootScope.$emit('location-saved', false);
	        }).error(function(data, status, headers) {
	            notifyError('Localisation', 'Erreur survenue lors de la localisation', data);
	        });
		},
		updateState: function(announcementId, action) {
			notifyInfo('Modification', 'Prise en compte de la demande de modification...');
			var deferred = $q.defer();
			$http.put('/ws/animaux/id/' + announcementId + '?action=' + action).success(function(data, status, headers) {
				deferred.resolve();
				notifySuccess('Modification', 'Modification effectuée avec succès');
			}).error(function(data, status, headers) {
				if (status == 403) {
					notifyError('Modification', 'Vous n\'avez pas les droits sur cette annonce');		
				} else {
					notifyError('Modification', 'Erreur survenue lors de la modification');		
				}
				deferred.reject();
			});
			return deferred.promise;
		},
		authenticate: function(announcementId, account) {
			$http.post('/login?id=' + announcementId, account).success(function(data, status, headers) {
				if (status == 200) {
					$rootScope.$emit('isAuthenticated', headers('Authorization'));
				}
				$rootScope.$emit('loading', false);
			}).error(function(data, status, headers) {
				if (status == 401) {
					notifyWarning('Authentification', 'Email / mot de passe incorrect pour cette annonce');
				} else {
					notifyWarning('Authentification', 'Erreur lors de l\'authentication.<br/> Veuillez réessayer, si cela persiste merci de contacter l\'administrateur');
				}
				$rootScope.$emit('loading', false);
			});
		}
	
	};
}]
);

app.service('dateService', function() {
	return {
		computeDuration: function(date) {
			return moment(date).locale('fr').calendar();
		},
		getFormattedDate: function(date) {
			return moment(date).locale('fr').fromNow();
		},		
		fetchLastLocationDate: function(locations) {
			if (locations && locations.length > 0) {
				return moment(locations[locations.length - 1].Date).locale('fr').calendar();
			}
		}
	};
});

app.service('utils', function() {
	return {
		getPicture: function(picture) {
			return picture ? '/photos/' + picture : '/static/images/noimage.gif';
		},
		displayPhoneNumber: function(announcement) {
			announcement.displayedPhoneNumber = announcement.PhoneNumber;
		}		
	};
});


app.controller('IndexCtrl', ['$scope', '$http', 'announcementService', 'utils', function($scope, $http, announcementService, utils) {	
	
	$scope.getOneAnnouncement = function() {
		announcementService.getAnnouncements(0, 'perdu').then(function(result) {
			$scope.announcements = result.announcements;
			$scope.totalItems = result.totalItems;
		});
	};
		
	$scope.searchCriteria = function() {
		if ($scope.criteria && $scope.criteria != '') {
			window.location = '/animaux?criteria=' + $scope.criteria;		
		} else {
			window.location = '/animaux';		
		}
	};
	
	$scope.getThumbnail = function(value) {
		return getPicture(value.img);
	};
	
	function getPicture(picture) {
		return picture ? picture : '/static/images/noimage.gif';
	}
	
	$scope.searchSuggestion = function(value) {
		var suggestions = [];
		var indexFound = 0;
		for (var i = 0;i < $scope.cities.length && i < 5;i++) {
			var cityLabel = $scope.cities[i].label;
			if (cityLabel.toLowerCase().indexOf(value) > -1) {
				suggestions.push({img: '/static/images/gps.png', label: cityLabel, category: 'city', index: ++indexFound});
			}
		}
		return $http.get('/ws/animaux?offset=0&limit=5&criteria=' + value)
		.then(function(response){
			if (response.data.announcements) {
				var indexFound = 0;
				for (var i = 0;i < response.data.announcements.length;i++) {
					var announcement = response.data.announcements[i];
					suggestions.push({img: '/photos/' + announcement.Picture, label: announcement.Name, category: 'announcement', index: ++indexFound});
				}
			}
			return suggestions;
		});
		return suggestions;
	};
	
	$scope.getOneAnnouncement();
}]
);

app.controller('ContactCtrl', ['$scope', '$http', function($scope, $http) {	
	
	$scope.sendMail = function() {
		if ($scope.contactForm.$invalid) {
			$scope.contactFormSubmitted = true;
			return;
		}
		$scope.contactFormSubmitted = false;
		$http.post('/ws/mail', { sender: $scope.mail.Sender, message: $scope.mail.Message, announcementId: $scope.selectedAnnouncement.Id }).success(function(data, status, headers) {
			$('#contact').modal('hide');
			$scope.mail = undefined;			
			notifySuccess('Contact', 'Message envoyé');
		}).error(function(data, status, headers) {
			notifyError('Contact', 'Erreur lors l\'envoi du message.<br/> Veuillez réessayer, si cela persiste merci de contacter l\'administrateur');
		});
	};
}]
);

app.controller('OneAnnouncementCtrl', ['$scope', '$http', 'announcementService', 'utils', 'dateService', function($scope, $http, announcementService, utils, dateService) {	

	$scope.searchById = function(id) {
		$http.get('/ws/animaux?id=' + id + '&offset=0&limit=1').success(function(data, status, headers) {
			$scope.announcements = data.announcements;
		}).error(function(data, status, headers) {
			notifyError('Recherche', 'Erreur lors de l\'exécution de la recherche.<br/> Veuillez réessayer, si cela persiste merci de contacter l\'administrateur');
		});
	};
	
	$scope.addLocation = function() {
		announcementService.addLocation($scope.selectedAnnouncement, $scope.location);
	};

	$scope.prepareUpdate = function(action, announcement) {
		$scope.actionRequiredAuth = action;
		$scope.announcementToUpdate = announcement;
	};	
		
	$('#authentication').on('hide.bs.modal', function (e) {
		$scope.account = undefined;
	});
	$('#authentication').on('shown.bs.modal', function () {
    	$('#email').focus()
  	});
	$('#contact').on('shown.bs.modal', function () {
    	$('#sender').focus()
  	});
	$('#location').on('shown.bs.modal', function () {
		if ($scope.selectedAnnouncement.Locations && $scope.selectedAnnouncement.Locations.length > 0) {
			var lastLocation = $scope.selectedAnnouncement.Locations[$scope.selectedAnnouncement.Locations.length - 1];		
			initializeMaps(true, $scope, 12, { lat: lastLocation.Latitude, lng: lastLocation.Longitude}, $scope.selectedAnnouncement.Locations);
		} else {
			initializeMaps(true, $scope, 10, { lat: -21.1306889, lng: 55.5264794});
		}		
	});
	
	var url = window.location.pathname.split('/');
	if (url.length >= 4) {
		$scope.searchById(url[3]);
	} else {
		notifyError('Animal perdu et errant', 'Le paramètre identifiant est manquant');
	}
	$scope.dateService = dateService;
}]
);
	
app.controller('LostListCtrl', ['$rootScope', '$scope', '$http', 'announcementService', 'dateService', 'utils', function($rootScope, $scope, $http, announcementService, dateService, utils) {
	
	$scope.getAnnouncements = function(offset, type, searchQueryParam) {
		if (offset == 0) {
			$scope.currentPage = 1;
		}
		announcementService.getAnnouncements(offset, type, searchQueryParam).then(function(result) {
			$scope.announcements = result.announcements;
			$scope.totalItems = result.totalItems;
		});
	};
	
	$scope.getAnnouncementsFromType = function(type) {
		if (type != $scope.type) {
			$scope.type = type;
			$scope.criteria = undefined;
			$scope.city = undefined;
			$scope.getAnnouncements(0, type);
		}
	};	
	
	$rootScope.$on('loading', function(event, loading) {
		event.stopPropagation();
		$scope.displayLoading(loading);
	});
	
	$scope.searchCity = function(offset) {
		$scope.search('city', offset);
	};
	
	$scope.searchCriteria = function(offset) {
		$scope.search('criteria', offset);
	};
	
	$scope.search = function(queryParam, offset) {
		offset = offset ? offset : 0;
		$scope.find(queryParam, $scope.criteria, offset);
	};		
	
	$scope.addLocation = function() {
		announcementService.addLocation($scope.selectedAnnouncement, $scope.location);
	};
	
	$scope.prepareUpdate = function(action, announcement) {
		$scope.actionRequiredAuth = action;
		$scope.announcementToUpdate = announcement;
	};
	
	$('#location').on('shown.bs.modal', function () {
		if ($scope.selectedAnnouncement.Locations && $scope.selectedAnnouncement.Locations.length > 0) {
			var lastLocation = $scope.selectedAnnouncement.Locations[$scope.selectedAnnouncement.Locations.length - 1];		
			initializeMaps(true, $scope, 12, { lat: lastLocation.Latitude, lng: lastLocation.Longitude}, $scope.selectedAnnouncement.Locations);
		} else {
			initializeMaps(true, $scope, 10, { lat: -21.1306889, lng: 55.5264794});
		}		
	});
		
	$scope.find = function(searchQueryParam, criteria, offset) {
		var searchQueryParam = searchQueryParam + '=' + $scope.criteria;
		$scope.getAnnouncements(offset, undefined, searchQueryParam);
	};
	
	$scope.pageChanged = function() {
		var offset = ($scope.currentPage - 1) * $rootScope.numberByPage;
		$scope.getAnnouncements(offset, $scope.type, $scope.querySearchParam);
		var target_offset = $('.search-result').offset();
        var target_top = target_offset.top;
        $('html, body').animate({scrollTop:target_top}, 500);
	};
	
	$scope.querySearchParam = undefined;
	var url = decodeURI(window.location.toString()).split('?');
	if (url.length >= 2) {
		var splittedQueryParam = /^(id|criteria|city)=(.+)/g.exec(url[1]);
		
		if (splittedQueryParam[1] == 'id') {
			$scope.searchById(splittedQueryParam[2]);
			return;
		}
		$scope.querySearchParam = splittedQueryParam[0];
		$scope.criteria = splittedQueryParam[2];
	}
	
	$scope.getAnnouncements(0, $scope.type, $scope.querySearchParam);
	$scope.dateService = dateService;
	$scope.currentPage = 1;
	$scope.maxSize = 3;
	$scope.numPages = 15;
	
}]
);

app.controller('AuthCtrl', ['$rootScope', '$scope', '$http', '$q', 'announcementService', function($rootScope, $scope, $http, $q, announcementService) {
		
	$scope.launchUpdate = function(action, announcement) {
		closeAuth();
		var deferred = $q.defer();
		deferred.promise.then(announcementService.updateState(announcement.Id, action))
		.then(announcementService.getAnnouncements(0, $scope.type).then(function(result) {
			$scope.$parent.announcements = result.announcements;
			$scope.$parent.totalItems = result.totalItems;
		}));
	};
	
	$scope.authenticate = function() {
		announcementService.authenticate($scope.announcementToUpdate.Id, $scope.account);
	};			
		
	$rootScope.$on('isAuthenticated', function(event, token, action) {
		$http.defaults.headers.common['Authorization'] = token;
		$scope.launchUpdate($scope.actionRequiredAuth, $scope.announcementToUpdate);
	});	
}]);

app.controller('AdminCtrl', ['$scope', '$http', function($scope, $http) {
	$scope.authenticate = function() {
		if ($scope.loginForm.$invalid) {
			$scope.submitted = true;
			return;
		}
		$http.post('/admin/login', $scope.account).success(function(data, status, headers) {
			if (status == 200) {
				$scope.$emit('isAuthenticated', headers('Authorization'));
			}			
		}).error(function(data, status, headers) {
			if (status == 401) {
				notifyWarning('Authentification', 'Email / mot de passe incorrect pour cette annonce');
			} else {
				notifyWarning('Authentification', 'Erreur lors de l\'authentication.<br/> Veuillez réessayer, si cela persiste merci de contacter l\'administrateur');
			}
		});
	};
	
	$scope.$on('isAuthenticated', function(event, data) {
		$http.defaults.headers.common['Authorization'] = data;
		fetchAddAnnouncementsTemplate();
	});
	
	function fetchAddAnnouncementsTemplate() {
		$http.get('/admin/announcements').success(function(data, status) {
			$scope.addAnnouncementsTemplate = data;
		})
		.error(function(status, data) {
			notifyError('Annonces', 'Erreur lors de la récupération du contenu des annonces');
		});
	}
	
	$scope.availableStates = ['waiting for validation', 'validated', 'deleted', 'found', 'deactivated'];
	
	$scope.getAllAnnouncements = function() {
		$scope.displayLoading(true);
		$http.get('/ws/admin/announcements/lost/all?state=' + $scope.state).success(function(data, status, headers) {
			$scope.announcements = data;
			$scope.displayLoading(false);
		}).error(function(data, status, headers) {
			notifyError('Liste des animaux perdus', 'Erreur lors de la récupération de la liste des animaux perdus.<br/> Veuillez réessayer, si cela persiste merci de contacter l\'administrateur');
			$scope.displayLoading(false);
		});
	};
	$scope.updateState = function(announcement, state) {
		notifyInfo('Modification', 'Prise en compte de la demande de modification');
		$http.put('/ws/admin/announcements/id/' + announcement.Id + '?state=' + state).success(function(data, status, headers) {
			notifySuccess('Modification', 'Modification effectuée');
			$scope.getAllAnnouncements();
		}).error(function(data, status, headers) {
			notifyError('Modification', 'Erreur lors de la modification');
		});
	};
}]);

app.controller('CreateCtrl', ['$scope', '$http', function($scope, $http) {
	$scope.maxDate = new Date();	
	
	$scope.isDateInvalid = function() {
		if ($scope.announcement == undefined || $scope.announcement.LostDate == undefined) {
			return false;
		}
		return !moment($scope.announcement.LostDate, momentDateFormat).locale('fr').isValid();
	};
	
	$scope.geoLocationActivated = function() {
		return navigator.geolocation;
	};
	
	$scope.geoLocalize = function() {
		if(navigator.geolocation) {
			navigator.geolocation.getCurrentPosition(function(position) {
				var pos = new google.maps.LatLng(position.coords.latitude, position.coords.longitude);
				map.setCenter(pos);
				map.setZoom(13);
				var locationMarker;
				var image = {
					url: '/static/images/location.png',
					size: new google.maps.Size(45, 45),
					origin: new google.maps.Point(0, 0),
					anchor: new google.maps.Point(17, 34),
					scaledSize: new google.maps.Size(15, 20)
				};
				locationMarker = new google.maps.Marker({ position: pos, map: map, animation: google.maps.Animation.DROP, icon: image });
				window.setTimeout(function() {
					locationMarker.setMap(null);
				}, 3000);
			});
		}
	};
	
	$scope.isPasswordsDifferent = function() {
		return $scope.announcement && $scope.announcement.Account && $scope.announcement.Account.Password != $scope.announcement.Account.Confirmation;
	};
	
	$scope.groupByRegion = function (item){
		return angular.uppercase(item.region);
	};
	
	$scope.searchLocation = function () {
		var geocodeRequest = {
			'address': $scope.locationAddress,
			'latLng': new google.maps.LatLng(-21.1306889, 55.5264794),
			'region': 'fr',
			'componentRestrictions': {'country': 'Réunion'}
		};
		var geocoder = new google.maps.Geocoder();
		geocoder.geocode(geocodeRequest, function(results, status) {
			if (status == google.maps.GeocoderStatus.OK) {
				var firstAddress = results[0];
				map.panTo(firstAddress.geometry.location);
				map.setZoom(15);
			}
		});
	};
	
	$scope.createAnnouncement = function(type) {
		if ($scope.createForm.$invalid || $scope.isPasswordsDifferent()) {
			$scope.submitted = true;
			return;
		}
		$scope.announcement.Type = type;
		$scope.announcement.Locations = [ $scope.location ];
				
	    notifyInfo('Signalement', 'Prise en compte de la demande de signalement d\'un animal...')
		$http.post('/ws/animaux/' + $scope.announcement.Type, $scope.announcement)
		.success(function(data, status, headers) {
            notifySuccess('Signalement', 'Enregistrement du signalement de l\'animal effectué avec succès')
	        $scope.announcement = undefined;
			$scope.submitted = false;
	        $scope.displayLoading(false);
        }).error(function(data, status, headers) {
            notifyError('Signalement', 'Erreur lors du signalement d\'un animal.<br/> Veuillez réessayer, si cela persiste merci de contacter l\'administrateur', data)
            $scope.displayLoading(false);
        });
	};
	
	$scope.open = function($event) {
	    $event.preventDefault();
	    $event.stopPropagation();
		
	    $scope.opened = true;
	};
		
	$scope.fetchSpecies = function() {
		$http.get('/ws/species').success(function(data, status, headers) {
			$scope.species = data;
		}).error(function(data, status, headers) {
			notifyError('Race', 'Erreur lors de la récupération de la liste des races de chien.<br/> Veuillez réessayer, si cela persiste merci de contacter l\'administrateur');
		});
	}();

	$scope.initMaps = function() {	
		initializeMaps(false, $scope, 10, { lat: -21.1306889, lng: 55.5264794});
	};
	
	$scope.submitted = false;
	$scope.announcement = {};
	$scope.dateFormat = 'EEEE dd MMMM yyyy';
	var momentDateFormat = 'dddd D MMMM YYYY';
}]
);

var map;

function initializeMaps(hasValidateBtn, $scope, zoom, position, locations) {
	var mapOptions = {
		center: position,
		zoom: zoom,
		panControl: false,
		zoomControl: true,
		zoomControlOptions: {
			style: google.maps.ZoomControlStyle.SMALL
		},
		mapTypeControl: false,
		scaleControl: false,
		streetViewControl: true,
		overviewMapControl: false
	};
	
	map = new google.maps.Map(document.getElementById('google-maps'), mapOptions);
	var infoDiv = document.createElement('div');
	var displayInfo = new DisplayInfo(infoDiv, map);
	infoDiv.index = 1;
	map.controls[google.maps.ControlPosition.TOP_LEFT].push(infoDiv);

	if (hasValidateBtn) {
		var validateDiv = document.createElement('div');
		var validateControl = new ValidateControl(validateDiv, map);
		validateDiv.index = 1;
		map.controls[google.maps.ControlPosition.RIGHT_BOTTOM].push(validateDiv);
		google.maps.event.addDomListener(validateDiv, 'click', function() {
			$scope.addLocation()
		});
	}
	
	var markers = [];		
	if (locations) {
		var lastLocation = false;
		var coordonates = [];
		var lastCoordonate;
		for (var index in locations) {
			if (index == locations.length - 1) {
				lastLocation = true;
			}
			var coordonate = new google.maps.LatLng(locations[index].Latitude, locations[index].Longitude);
			placeMarker(coordonate, locations[index].Date, lastLocation);
			coordonates.push(coordonate);
			if (lastCoordonate) {
				computeItinery(lastCoordonate, coordonate);
			}
			lastCoordonate = coordonate;			
		}
		var way = new google.maps.Polyline({
			path: coordonates,        
			strokeColor: "#FF0000",
			strokeOpacity: 0,
			strokeWeight: 1,
			icons: [{
		      icon: { path: 'M 0,-1 0,1', strokeWeight: 1, strokeOpacity: 0.5, scale: 4 },
		      offset: '0',
		      repeat: '20px'
		    }]
		});
		way.setMap(map);
	}
	
	function computeItinery(origin, destination) {
		var request = {
            origin      : origin,
            destination : destination,
            travelMode  : google.maps.DirectionsTravelMode.WALKING
        }
        var directionsService = new google.maps.DirectionsService();
        directionsService.route(request, function(response, status){
            if(status == google.maps.DirectionsStatus.OK){
                new google.maps.DirectionsRenderer({
				    map    : map,
					options : { suppressMarkers: true }
				}).setDirections(response);
            }
        });
	}
	
	google.maps.event.addListener(map, 'click', function(e) {
		clearMarkers();
		placeTempMarker(e.latLng, new Date());
	});
	
	function placeMarker(position, date, lastLocation) {
		var animal = $scope.selectedAnnouncement && $scope.selectedAnnouncement.Animal ? $scope.selectedAnnouncement.Animal : $scope.announcement.Animal;
		var image = {
			url: animal == 'Chat' ? '/static/images/cat-maps.png' : '/static/images/dog-maps.png',
			size: new google.maps.Size(45, 45),
			origin: new google.maps.Point(0, 0),
			anchor: new google.maps.Point(17, 34),
			scaledSize: new google.maps.Size(45, 45)
		};
		var marker = new google.maps.Marker({
			position: position,
			map: map,
			animation: google.maps.Animation.DROP,
			icon: image
		});
		
		function createMarker(address) {
			var markerinfo = new google.maps.InfoWindow({
				content: (lastLocation ? 'Dernière localisation<br/>' + moment(date).locale('fr').calendar() : moment(date).locale('fr').calendar()) + (address ? ' <br/>Adresse : ' + address : '')
			});
			if (lastLocation) {
				markerinfo.open(map, marker);
			}
	
			google.maps.event.addListener(marker, 'click', function() {
				markerinfo.open(map, marker);
			});
		}
		
		new google.maps.Geocoder().geocode({location: position}, function(results, status) {
			if (status == google.maps.GeocoderStatus.OK && results.length > 0) {
				createMarker(results[0].formatted_address);
			}
		});
		
		return marker;
	}

	function placeTempMarker(position, date) {
		var marker = placeMarker(position, date);
		markers.push(marker);
		$scope.location = { Latitude: marker.position.lat(), Longitude: marker.position.lng(), Date: date};
	}
	
	function clearMarkers() {
		for (var i = 0; i < markers.length; i++) {
			markers[i].setMap(null);
		}
		markers = [];
	}
}


function DisplayInfo(infoDiv, map) {
	infoDiv.style.padding = '5px';
	
	var infoUI = document.createElement('div');
	infoUI.style.backgroundColor = 'white';
	infoUI.style.borderStyle = 'solid';
	infoUI.style.borderWidth = '2px';
	infoUI.style.textAlign = 'center';
	infoDiv.appendChild(infoUI);
	
	var infoText = document.createElement('div');
	infoText.style.fontFamily = 'Arial,sans-serif';
	infoText.style.fontSize = '12px';
	infoText.style.paddingLeft = '4px';
	infoText.style.paddingRight = '4px';
	infoText.innerHTML = 'Faites un <b>clic</b> sur la carte <br/>pour localiser l\'animal';
	infoUI.appendChild(infoText);
}


function ValidateControl(validateDiv, map) {
	validateDiv.style.padding = '5px';
	
	var validateUI = document.createElement('div');
	validateUI.style.backgroundColor = 'white';
	validateUI.style.borderStyle = 'solid';
	validateUI.style.borderWidth = '2px';
	validateUI.style.cursor = 'pointer';
	validateUI.style.textAlign = 'center';
	validateDiv.appendChild(validateUI);
	
	var validateText = document.createElement('div');
	validateText.style.fontFamily = 'Arial,sans-serif';
	validateText.style.fontSize = '15px';
	validateText.style.paddingLeft = '14px';
	validateText.style.paddingRight = '14px';
	validateText.innerHTML = 'Valider';
	validateUI.appendChild(validateText);
}


function closeAuth() {
	$('#authentication').modal('hide');
}

app.directive('restrictedContent', function($compile, $parse) {
    return {
      restrict: 'E',
      link: function(scope, element, attr) {
        scope.$watch(attr.content, function() {
        		element.html($parse(attr.content)(scope));
        		$compile(element.contents())(scope);
        }, true);
      }
    }
});

app.directive("fileread", [function () {
    return {
        scope: {
            fileread: "="
        },
        link: function (scope, element, attributes) {
            element.bind("change", function (changeEvent) {
                var reader = new FileReader();
                reader.onload = function (loadEvent) {
                    scope.$apply(function () {
                        scope.fileread = loadEvent.target.result;
                    });
                }
				if (changeEvent.target.files[0] != undefined) {
                	reader.readAsDataURL(changeEvent.target.files[0]);
				}
            });
        }
    }
}]);
