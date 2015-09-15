var data = [];

var NameTagsModel = function (nameTags) {
    this.nameTags = ko.observableArray(ko.utils.arrayMap(nameTags, function (nameTag) {
        return {
            Name: nameTag.Name
        };
    }));
};

var ntm = new NameTagsModel(data);

ko.applyBindings(ntm);

//ntm.nameTags([{Name:"Tester"}]);

function getData() {
    $.ajax({
        url: '/queue/getAll',
        dataType: 'json',
        success: function(data) {
            console.log(data);
            ntm.nameTags(data);
        }
    })
}

$(document).ready(function () {
    getData();
    window.setInterval(getData, 1000);
});