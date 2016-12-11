;;; gotn.el --- Utility tool for Go programming language to run test at point

;; Author: Yauhen Lazurkin
;; Keywords: languages go test
;; URL: https://github.com/yauhenl/gotn.el

;;; Code:

(defcustom go-test-case-command "go test -v -run"
  "The 'godef' command."
  :type 'string
  :group 'gotn)

(defun gotn-run-test (point)
  "Run go test at POINT."
  (interactive "d")
  (condition-case nil
      (let ((gotn-out (gotn--call point)))
        (if (= (car gotn-out) 0)
          (shell-command
           (concat go-test-case-command " ^" (car (cdr gotn-out)) "$"))
          (message (car (cdr gotn-out)))))
    (file-error (message "Could not run gotn binary"))))

(defun gotn--call (point)
  "Call `gotn' to get test name at POINT."
  (if (not (buffer-file-name (current-buffer)))
      (error "Cannot use gotn on a buffer without a file name")
    (let ((out
           (shell-command-to-string
            (concat "gotn -f "
                    (file-truename (buffer-file-name (current-buffer)))
                    " -p "
                    (number-to-string (position-bytes point))))))
      (if (string= (substring out -1 nil) "\n")
          (list 1 (substring out 0 -1))
        (list 0 out)))))

(provide 'gotn)

;;; gotn.el ends here
