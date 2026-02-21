$(document).ready(function () {
    $('.categories-selection-multiple').select2();
})



function filterPosts() {
    //Get the select select list and store in a variable
    var filterCategories = document.getElementById("filter_categories");
    //Get the selected value of the select list
    var formaction = filterCategories.options[filterCategories.selectedIndex].value;

    document.posts_filter.action = window.origin + '/categories/' + formaction;
}