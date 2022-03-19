$(document).ready(function () {

    $('form.j-validate').validate({
        submitHandler: function (form) {
            let url = $(form).attr('action') || location.href;
            $.post(url, $(form).serialize(), function (ret) {
                if (ret.code > 0) {
                    layer.msg(ret.msg)
                    return;
                }
                layer.msg('success')
            }, 'json');
        }
    });

});