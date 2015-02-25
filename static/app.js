(function() {
  var map = L.map('map');
  L.tileLayer.provider('OpenStreetMap.BlackAndWhite').addTo(map);
  map.locate({setView: true, maxZoom: 16});
})();
