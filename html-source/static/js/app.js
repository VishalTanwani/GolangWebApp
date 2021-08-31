let myFunc = Prompt();

function Prompt() {
    let toast = function (obj) {
        const { msg = "", icon = "success", position = "top-end" } = obj;

        const Toast = Swal.mixin({
            toast: true,
            position: position,
            icon: icon,
            title: msg,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
            toast.addEventListener("mouseenter", Swal.stopTimer);
            toast.addEventListener("mouseleave", Swal.resumeTimer);
            },
        });

        Toast.fire({});
    };

    let success = function (obj) {
        const { msg = "", title = "", footer = "" } = obj;

        Swal.fire({
            icon: "success",
            title: title,
            text: msg,
            footer: footer,
        });
    };

    let error = function (obj) {
        const { msg = "", title = "", footer = "" } = obj;

        Swal.fire({
            icon: "error",
            title: title,
            text: msg,
            footer: footer,
        });
    };

    let custom = async function (obj) {
        const { title = "", html = "", showConfirmButton = true } = obj;
        const { value: result } = await Swal.fire({
            icon:obj.icon,
            title: title,
            html: html,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            focusConfirm: false,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if(obj.willOpen !== undefined) {
                    obj.willOpen()
                }
            },
            didOpen: () => {
                if(obj.didOpen !== undefined) {
                    obj.didOpen()
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById("start_date").value,
                    document.getElementById("end_date").value,
                ];
            },
        });

        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    if (obj.callback !== undefined) {
                        obj.callback(result)
                    }
                } else {
                    obj.callback(false)
                }
            } else {
                obj.callback(false)
            }
        }
    };

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    };
}