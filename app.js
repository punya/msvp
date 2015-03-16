import {Incident} from './incident';

var map = L.map('map');
L.tileLayer('http://otile{s}.mqcdn.com/tiles/1.0.0/map/{z}/{x}/{y}.jpeg', {
  attribution: `
    Tiles Courtesy of <a href="http://www.mapquest.com/">MapQuest</a>
    &mdash; Map data &copy;
    <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>
  `,
  subdomains: '1234'
}).addTo(map);
var markers = new L.MarkerClusterGroup();
map.addLayer(markers);
map.locate({
  setView: true,
  maxZoom: 16
});
map.on('contextmenu', (e) => {
  $("#location").val(e.latlng.lat + "," + e.latlng.lng);
  $("#reportModal").modal();
});

$.getJSON('incidents', (data) => {
  for (let json of data) {
    new Incident(json.Text, json.Lat, json.Lng).render(markers);
  }
});

$('#submit').click(() => {
  let loc = $('#location').val().split(',');
  let lat = NaN, lng = NaN;
  if (loc.length == 2) {
    lat = parseFloat(loc[0]);
    lng = parseFloat(loc[1]);
  }
  let text = $('textarea').val();
  $.post('incidents', JSON.stringify({lat, lng, text}));
});
