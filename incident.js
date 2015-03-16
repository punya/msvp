export class Incident {
  constructor(text, lat, lng) {
    this.text = text;
    this.lat = lat;
    this.lng = lng;
  }

  render(markers) {
    var marker = L.marker([this.lat, this.lng]);
    marker.bindPopup(this.text);
    markers.addLayer(marker);
  }

  save(callback) {
    var req = new XMLHttpRequest();
    req.onload = callback;
    req.open('post', 'incidents', true);
    req.setRequestHeader('Content-Type', 'application/json');
    req.send(JSON.stringify(this));
  }
}
