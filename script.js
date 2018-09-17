// var x = document.getElementById("demo");
var markers = [];
function getLocation() {
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(showPosition);
    } else {
        x.innerHTML = "Geolocation is not supported by this browser.";
    }
}

function showPosition(position) {
    document.getElementById("current_latitude").value= position.coords.latitude;
    document.getElementById("current_longitude").value= position.coords.longitude;
    initMap();
}

function initMap() {
var lat =  parseFloat(document.getElementById("current_latitude").value);
var lng = parseFloat(document.getElementById("current_longitude").value);
var map = new google.maps.Map(document.getElementById('map'), {
  mapTypeControl: false,
  center: {lat: lat, lng: lng},
  zoom: 13
});

var geocoder = new google.maps.Geocoder;
var infowindow = new google.maps.InfoWindow({disableAutoPan:true});
geocodeLatLng(geocoder, map, infowindow);
}

function geocodeLatLng(geocoder, map, infowindow) {
     var lat =   parseFloat(document.getElementById("current_latitude").value);
     var lng = parseFloat(document.getElementById("current_longitude").value);
     var latlng = {lat: lat, lng: lng};
     geocoder.geocode({'location': latlng}, function(results, status) {
       if (status === 'OK') {
         if (results[0]) {
           map.setZoom(11);
           var marker = new google.maps.Marker({
             position: latlng,
             map: map,
             draggable : true
           });
           google.maps.event.addListener(marker,'dragend',function(event) {
                var newlat = event.latLng.lat();
                var newlng = event.latLng.lng();
                var latlng = {lat: newlat, lng: newlng};
                geocoder.geocode({'location': latlng}, function(results, status) {
                    if (status === 'OK') {
                        document.getElementById("origin-input").value = results[0].formatted_address;
                        document.getElementById("orig_latitude").value= results[0].geometry.location.lat();
                        document.getElementById("orig_longitude").value= results[0].geometry.location.lng();
                    }
                });
            });
           markers.push(marker);
           document.getElementById("origin-input").value = results[0].formatted_address;
           infowindow.setContent(results[0].formatted_address);
           infowindow.open(map, marker);
           infowindow.close();
           document.getElementById("orig_latitude").value= results[0].geometry.location.lat();
           document.getElementById("orig_longitude").value= results[0].geometry.location.lng();
           var place_id = results[0].place_id;
           new AutocompleteDirectionsHandler(map, place_id);
         } else {
           window.alert('No results found');
           new AutocompleteDirectionsHandler(map);
         }
       } else {
         window.alert('Geocoder failed due to: ' + status);
       }
     });
   }

/**
* @constructor
*/
function AutocompleteDirectionsHandler(map, place_id) {
this.map = map;
this.originPlaceId = null;
this.destinationPlaceId = null;
this.travelMode = 'DRIVING';
var originInput = document.getElementById('origin-input');
var destinationInput = document.getElementById('destination-input');
this.originPlaceId = place_id;
//        var modeSelector = document.getElementById('mode-selector');
this.directionsService = new google.maps.DirectionsService;
this.directionsDisplay = new google.maps.DirectionsRenderer;
this.directionsDisplay.setMap(map);

var originAutocomplete = new google.maps.places.Autocomplete(
    originInput, {placeIdOnly: false});
var destinationAutocomplete = new google.maps.places.Autocomplete(
    destinationInput, {placeIdOnly: false});

this.setupPlaceChangedListener(originAutocomplete, 'ORIG');
this.setupPlaceChangedListener(destinationAutocomplete, 'DEST');
}

AutocompleteDirectionsHandler.prototype.setupPlaceChangedListener = function(autocomplete, mode) {
var me = this;
autocomplete.bindTo('bounds', this.map);
autocomplete.addListener('place_changed', function() {
  var place = autocomplete.getPlace();
  if (!place.place_id) {
    window.alert("Please select an option from the dropdown list.");
    return;
  }
  if (mode === 'ORIG') {
    me.originPlaceId = place.place_id;
    document.getElementById("orig_latitude").value= place.geometry.location.lat();
    document.getElementById("orig_longitude").value= place.geometry.location.lng();
  } else {
    me.destinationPlaceId = place.place_id;
    document.getElementById("dest_latitude").value= place.geometry.location.lat();
    document.getElementById("dest_longitude").value= place.geometry.location.lng();
  }
  me.route();
});
};

AutocompleteDirectionsHandler.prototype.route = function() {
for (var i = 0; i < markers.length; i++) {
    markers[i].setMap(null);
  }
if (!this.originPlaceId || !this.destinationPlaceId) {
  return;
}
var me = this;

this.directionsService.route({
  origin: {placeId: this.originPlaceId},
  destination: {placeId: this.destinationPlaceId},
  travelMode: this.travelMode
}, function(response, status) {
  if (status === 'OK') {
    me.directionsDisplay.setDirections(response);
  } else {
    window.alert('Directions request failed due to ' + status);
  }
});
};

$(function() {
  $("#search_form").submit(function(e){
       e.preventDefault();
   });
   $(".map_trigger").click(function(){
     $('#map').slideToggle("slow");
     $("#search_data, .context_toggle").toggleClass("visibility_hidden");
   });
   //Disable and enable the search button validation.
   $('#search_compare').prop('disabled',true);
       $("#origin-input, #destination-input").keyup(function(){
         $('#search_compare').prop('disabled', this.value == "" ? true : false);
     });
  $("#search_compare").on("click", function() {
       $('#map').hide("slow");
       $("#search_data, .context_toggle").removeClass("visibility_hidden");

       let a = $("#orig_latitude").val(),
           b = $("#orig_longitude").val(),
           c = $("#dest_latitude").val(),
           d = $("#dest_longitude").val();

       $('#loadingmessage').show();
       $("table").find("tr:gt(0)").remove(); // Refresh the contents of table after each search and compare.
       //Ajax call to render the Cloudless API and render the dynamic html
       $.ajax({
           type: 'get',
           url: 'https://1ei6bnqkm5.execute-api.us-east-1.amazonaws.com/dev4/cabBooking?start_latitude='+a+'&start_longitude='+b+'&end_latitude='+c+'&end_longitude='+d,
           crossDomain: true,
           contentType: "application/json",
           dataType: 'json',
           headers: {
             "x-api-key":"DEKU4CGInq4VaEu29wmqp8wuwUcFXUVhaNqGIQ1b"
           },
           success: function (data, textStatus, xhr) {
             if (data.cabs)
             {
                let json = data.cabs,
                    tr,
                    i = 0;
                for (; i < json.length; i++) {
                  tr = $('<tr/>');
                  tr.append("<td>" + json[i].company + "</td>");
                  tr.append("<td>" + json[i].cab + "</td>");
                  tr.append("<td>" + "$ " + json[i].Estimate + "</td>");
                  tr.append("<td>" + "<p>" +json[i].arriving+ "</p>" + "<button class='book_your_cab'>"+ "<span>" + " BOOK "+ "</span>" + "</button>" +"</td>");
                  $('table').first().append(tr);
                }
              } else {
                  $('#errorsContainer').html(data.message );
                  $('#map').show("slow");
              }
              $('#loadingmessage').hide();
           },
           error: function (xhr, textStatus, errorThrown) {
                $('.data_table th').hide();
                $('#errorsContainer').show();
                $('#errorsContainer').html(xhr.responseJSON.message);
                $('#loadingmessage').hide();
                }
       });
       //Add click on dynamically created button and redirect the window to the specified location.
       $(document).on('click','.book_your_cab',function(){
          window.open('https://www.uber.com/en-IN/cities/pune/','_blank');
       });
  });
});
