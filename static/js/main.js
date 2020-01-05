$(document).ready(function () {
    // update number of link
    updateNumberOfListElements();

    // add handlers for buttons
    $('body').on('click', '.remove-link', function (event) {
        var $target = $(event.target);
        $target.closest('.link-group').remove();
        updateNumberOfListElements();
    });

    $('.link-form').on('submit', function (event) {
        event.preventDefault();

        var result = [];
        $('.links .link-group').each((i, element) => {
            var url = $(element).find('.url').val();
            var pinned = $(element).find('.pinned').is(':checked');
            result.push({url: url, pinned: pinned});
        });

        $.ajax("/save", {
                data: JSON.stringify(result),
                contentType: 'application/json',
                type: 'POST'
            }
        ).done(function () {
            $('.link-form input[type=submit]').after(
                `<div class="save-alert alert alert-success" role="alert">Links saved successfully</div>`
            );
            console.info("Links saved")
        }).fail(function () {
            $('.link-form input[type=submit]').after(
                `<div class="save-alert alert alert-danger" role="alert">Error occured. Please try again in some time.</div>`
            );
            console.info("Links cannot be saved")
        }).always(function () {
            $('.alert').fadeOut(1000, function () {
                $(this).remove();
            })
        });
        return false;
    });

    $('.add-link').on('click', function (event) {
        var lastElem = $(".links tbody tr").length;
        $(".links tbody").append(
            `
            <tr class="d-flex link-group">
                <td  class="col-1" scope="row">${lastElem + 1}</td>
                <td  class="col-9"><input type="url" class="form-control url" required></td>
                <td  class="col-1"><input type="checkbox" class="form-control pinned"></td>
                <td  class="col-1"><a class="btn btn-danger remove-link"><i class="fas fa-eraser"></i></a></td>
            </tr>
            `
        )
    });
});

function updateNumberOfListElements() {
    $(".links tbody td:first-child").each((i, element) => {
        $(element).html(i + 1);
    })
}
