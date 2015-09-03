function preview() {
    var name = $('#name').val();
    $.ajax({
        url: "/preview",
        data: {
            "name": name
        },
        statusCode: {
            400: function () {
                console.log("Invalid request made for preview image - 400");
                $.notify({
                    title: 'Name Tag Submission Failed: ',
                    message: 'Invalid request made for preview image - 400',
                    icon: 'glyphicon glyphicon-remove'
                },{
                    type: 'danger'
                });
            },
            404: function () {
                console.log("Request for preview image could not be made - 404");
                $.notify({
                    title: 'Name Tag Submission Failed: ',
                    message: 'Request for preview image could not be made - 404',
                    icon: 'glyphicon glyphicon-remove'
                },{
                    type: 'danger'
                });
            },
            500: function() {
                console.log("Internal server error occurred while processing request for preview image - 500");
                $.notify({
                    title: 'Name Tag Submission Failed: ',
                    message: 'Internal server error occurred while processing request for preview image - 500',
                    icon: 'glyphicon glyphicon-remove'
                },{
                    type: 'danger'
                });
            }
        },
        success: function (response) {
            console.log(response);
            if(response.Code === 0) {
                $('#preview-image').attr('src', response.Image);
                $.notify({
                    message: 'Name Tag Preview Successfully Generated',
                    icon: 'glyphicon glyphicon-ok'
                },{
                    type: 'success'
                });
            } else {
                console.log(response.Error);
            }
        }
    });
}

function submit() {
    var name = $('#name').val();
    $.ajax({url: "queue/add",
        type: "POST",
        data: {
            "name": name
        },
        statusCode: {
            400: function () {
                console.log("Invalid submit request made - 400");
                $.notify({
                    title: 'Name Tag Submission Failed: ',
                    message: 'Invalid submit request made - 400',
                    icon: 'glyphicon glyphicon-remove'
                },{
                    type: 'danger'
                });
            },
            404: function () {
                console.log("Submit request could not be made - 404");
                $.notify({
                    title: 'Name Tag Submission Failed: ',
                    message: 'Submit request could not be made - 404',
                    icon: 'glyphicon glyphicon-remove'
                },{
                    type: 'danger'
                });
            },
            500: function() {
                console.log("Internal server error occurred while processing submit request - 500");
                $.notify({
                    title: 'Name Tag Submission Failed: ',
                    message: 'Internal server error occurred while processing submit request - 500',
                    icon: 'glyphicon glyphicon-remove'
                },{
                    type: 'danger'
                });
            }
        },
        success: function (data) {
            console.log(data);
            $('#name').val("");
            $('#preview-image').attr('src', '../static/assets/blank.png');
            $.notify({
                message: 'Name Tag Successfully Submitted',
                icon: 'glyphicon glyphicon-ok'
            },{
                type: 'success'
            });
        }
    });
}

$(document).ready(function() {
    $('#preview').click(function() {
        preview();
    });

    $('#submit').click(function() {
        submit();
    });
});