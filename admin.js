$.getJSON('incidents?verified=false', function(incidents) {
    let content = $('.container');
    for (let {key, lat, lng, text, verified} of incidents) {
        let form = $(`
            <div class="panel panel-default">
                <div class="panel-body">
                    <a href="//maps.google.com?q=${lat},${lng}"><img src="//maps.googleapis.com/maps/api/staticmap?markers=${lat},${lng}&size=600x400&zoom=10"></a>
                    <br>
                    <div class="form-group">
                        <label>Latitude</label>
                        <input type="text" value="${lat}" class="form-control lat">
                    </div>
                    <div class="form-group">
                        <label>Longitude</label>
                        <input type="text" value="${lng}" class="form-control lng">
                    </div>
                    <div class="form-group">
                        <label>Text</label>
                        <textarea class="form-control text">${text}</textarea>
                    </div>
                </div>
                <div class="panel-footer">
                    <button class="btn btn-primary">Verify</button>
                    <button class="btn btn-danger">Delete</button>
                </div>
            </div>
        `);
        $('.btn-primary', form).click(() => {
            let data = JSON.stringify({
                lat: parseFloat($('.lat', form).val()),
                lng: parseFloat($('.lng', form).val()),
                text: $('textarea', form).val(),
                verified: true
            });
            $.ajax(`incidents/${key}`, {
                contentType: 'application/json',
                type: 'PUT',
                data: data
            }).then(() => {
                form.remove();
            });
        });
        $('.btn-danger', form).click(() => {
            $.ajax(`incidents/${key}`, {
                type: 'DELETE'
            }).then(() => {
                form.remove();
            });
        });
        form.appendTo(content);
    }
});
