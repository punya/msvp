(function() {
  "use strict";

  $("#about-btn").click(function() {
    $("#aboutModal").modal("show");
    $(".navbar-collapse.in").collapse("hide");
    return false;
  });

  var map = L.map('map');
  L.tileLayer('http://otile{s}.mqcdn.com/tiles/1.0.0/map/{z}/{x}/{y}.jpeg', {
    attribution: 'Tiles Courtesy of <a href="http://www.mapquest.com/">MapQuest</a> &mdash; Map data &copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
    subdomains: '1234'
  }).addTo(map);
  var markers = new L.MarkerClusterGroup();
  map.addLayer(markers);
  map.locate({setView: true, maxZoom: 16});

  function Incident(text, lat, lng) {
    this.text = text;
    this.lat = lat;
    this.lng = lng;
  }
  Incident.prototype.render = function render() {
    var marker = L.marker([this.lat, this.lng]);
    marker.bindPopup(this.text);
    markers.addLayer(marker);
  }
  Incident.prototype.save = function save(callback) {
    var req = new XMLHttpRequest();
    req.onload = function saveXhrCallback() {
      if (req.status == 200) {
        callback(JSON.parse(req.responseText));
      } else {
        console.log(req.responseText);
      }
    };
    req.open('post', 'incidents', true);
    req.setRequestHeader('Content-Type', 'application/json');
    req.send(JSON.stringify(this));
  }

  var req = new XMLHttpRequest();
  req.onload = function loadXhrCallback() {
    JSON.parse(req.responseText).forEach(function makeObject(json) {
      new Incident(json.Text, json.Lat, json.Lng).render();
    });
  };
  req.open('get', 'incidents', true);
  req.send();

  map.on('contextmenu', function handleLongPress(e) {
    var marker = L.marker(e.latlng);

    var container = document.createElement('div');
    container.innerHTML = '<textarea></textarea><br><button>Post</button><span></span>';
    container.querySelector('button').addEventListener('click', function handlePostButton() {
      var storyTextArea = container.querySelector('textarea');
      var savingMessage = container.querySelector('span');
      savingMessage.innerHTML = 'Saving&hellip;';
      var incident = new Incident(storyTextArea.value, e.latlng.lat, e.latlng.lng);
      incident.render();
      incident.save(function afterSave(e) {
        marker.closePopup();
      });
    });
    marker.on('popupclose', function handlePostPopupClose() {
      map.removeLayer(marker);
    });

    marker.addTo(map);
    marker.bindPopup(container);
    marker.openPopup();
  });
})();
