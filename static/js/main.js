$(document).ready(function () {
    //todo: make request for links using JWT
    var links = [
        {
            "link": "https://google.com",
            "pinned": true
        },
        {
            "link": "https://facebook.com",
            "pinned": false
        },
        {
            "link": "https://facebook.com",
            "pinned": false
        },
    ];
    // render links

    links.forEach(function(elem, index) {
        // todo: add validation for input
        $(".links tbody").append(
            `
            <tr class="link-group">
                <th scope="row">${index + 1}</th>
                <td><input type="url" class="form-control url" value="${elem.link}" required></td>
                <td><input type="checkbox" class="form-control pinned" value="${elem.pinned}"></td>
                <td><a class="btn btn-danger remove-link"><i class="fas fa-eraser"></i></a></td>
            </tr>
            `
        )
    });

    // add handlers for buttons
    $('body').on('click', '.remove-link', function (event) {
        var $target = $(event.target);
        $target.closest('.link-group').remove();
        $(".links tbody th:first-child").each((i, element) => {
            $(element).html(i + 1);
        })
    });

    $('.save-links').on('click', function (event) {
        var result = $('.links .link-group').map((i, element) => {
            var link = $(element).find('.url').val();
            var pinned = JSON.parse($(element).find('.pinned').val());
            return { link: link, pinned: pinned }
        });
        console.log(result);
        return false;
    });
});