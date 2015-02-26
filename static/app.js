(function() {
  var map = L.map('map');
  L.tileLayer.provider('OpenStreetMap').addTo(map);
  map.locate({setView: true, maxZoom: 16});

  function Incident(text, latitude, longitude) {
    this.text = text;
    this.latitude = latitude;
    this.longitude = longitude;
  }
  Incident.prototype.render = function render() {
    var marker = L.marker([this.latitude, this.longitude]).addTo(map);
    marker.bindPopup(this.text);
  }

  var incidents = [];
  incidents.forEach(function render(inc) { inc.render(); });

  map.on('contextmenu', function handleLongPress(e) {
    var marker = L.marker([e.latlng.lat, e.latlng.lng]);

    var container = document.createElement('div');
    container.innerHTML = '<textarea></textarea><br><button>Post</button>';
    var storyTextArea = container.querySelector('textarea');
    container.querySelector('button').addEventListener('click', function handlePostButton() {
      marker.closePopup();
      var incident = new Incident(storyTextArea.value, e.latlng.lat, e.latlng.lng);
      incidents.push(incident);
      incident.render();
    });
    marker.on('popupclose', function handlePostPopupClose() {
      map.removeLayer(marker);
    });

    marker.addTo(map);
    marker.bindPopup(container);
    marker.openPopup();
  });
})();
