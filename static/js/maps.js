var mapRaphael = angular.module('mapRaphael',[]);
mapRaphael.factory('mapFactory', function() {
    return map;
});

mapRaphael.directive('map', function ($rootScope, $window) {
    function createRapha(element) {
        var width = $('#map-wrapper').outerWidth(false);
        var paper = new Raphael(element, '100%', '100%');
        
        var attr = {
            fill: "#F4EAD6",
            stroke: "#CC5C3C",
            "stroke-width": 1,
            "stroke-linejoin": "round"
			
        };
        paper.canvas.setAttribute('viewBox', '62 59 978 900');
        paper.canvas.setAttribute('preserveAspectRatio', 'xMinYMin meet');
        paper.canvas.setAttribute('class', 'svg-content');
        $.get('/static/svg/reunion.svg', function(svg) {
			
            var map = {};
            $(svg).find('path').each(function (index, region) {
                var path = $(region).attr('d');
                var id =  $(region).attr('id');
                map[id] = paper.path(path).attr(attr);
            });
            
            $rootScope.$broadcast('map-data-ok', map);
        });
    }
    
    $rootScope.$on('map-data-ok', function (event, toDraw) {
        for (var state in  toDraw) {
            toDraw[state].color = "#000";
            (function (st, state) {
                st[0].id = state;
                st[0].style.cursor = "pointer";
                st[0].onmouseover = function () {
                    st.animate({fill: "#F9C866", stroke: "#CC5C3C"}, 300);
                    $rootScope.$broadcast('map-mouseover', st[0]);
                };
                
                st[0].onmouseout = function () {
                    st.animate({fill: "#F4EAD6", stroke: "#CC5C3C"}, 300);
                    $rootScope.$broadcast('map-mouseout', st[0]);
                };
      
                st[0].onmouseup = function () {
                    $rootScope.$broadcast('map-mouseup', st[0]);
                };
                
                st[0].onmousedown = function () {
                    $rootScope.$broadcast('map-mousedown', st[0]);
                };
            })( toDraw[state], state);
        }
    });
    
    return function (scope, element, attrs) {   
        createRapha(element[0]);
    };
});
