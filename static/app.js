(function() {
  "use strict";

  var map = L.map('map');
  L.tileLayer.provider('OpenStreetMap').addTo(map);
  map.locate({setView: true, maxZoom: 16});

  function Incident(text, lat, lng) {
    this.text = text;
    this.lat = lat;
    this.lng = lng;
  }
  Incident.prototype.render = function render() {
    var marker = L.marker([this.lat, this.lng]).addTo(map);
    marker.bindPopup(this.text);
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
    req.open('post', 'posts', true);
    req.setRequestHeader('Content-Type', 'application/json');
    req.send(JSON.stringify(this));
  }

  var req = new XMLHttpRequest();
  req.onload = function loadXhrCallback() {
    JSON.parse(req.responseText).forEach(function makeObject(json) {
      new Incident(json.Text, json.Lat, json.Lng).render();
    });
  };
  req.open('get', 'posts', true);
  req.send();

  map.on('contextmenu', function handleLongPress(e) {
    var marker = L.marker(e.latlng);

    var container = document.createElement('div');
    container.innerHTML = '<textarea></textarea><br><button>Post</button>';
    var storyTextArea = container.querySelector('textarea');
    container.querySelector('button').addEventListener('click', function handlePostButton() {
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
