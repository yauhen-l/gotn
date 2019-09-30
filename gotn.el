;;; gotn.el --- Utility tool for Go programming language to run test at point

;; Author: Yauhen Lazurkin
;; Keywords: languages go test
;; URL: https://github.com/yauhenl/gotn

;;; Commentary:

;; Quick start:
;; Install `gotn` tool with `go get github.com/yauhenl/gotn`
;; and make sure it is in the PATH
;;
;; Define a go-mode hook and bindings:
;; (require 'gotn)
;; (add-hook 'go-mode-hook
;;   (lambda ()
;;     (local-set-key (kbd "C-c t") 'gotn-run-test))
;;     (local-set-key (kbd "C-c C-t") 'gotn-run-test-package)))
;;
;; For more details visit: https://github.com/yauhenl/gotn

;;; Code:

(require 'compile)

(add-to-list 'compilation-error-regexp-alist '("\\([a-zA-Z0-9_]+\\.go\\):\\([0-9]+\\)" 1 2))

(defcustom go-test-case-command "go test -v -run"
  "The command to run test case."
  :type 'string
  :group 'gotn)

(defcustom go-test-package-command "go test -v ."
  "The command to run tests of current pacakge."
  :type 'string
  :group 'gotn)

(defcustom go-test-package-test-fallback t
  "Whether no test case under position run all package tests."
  :type 'boolean
  :group 'gotn)

(defun gotn--compilation-name (mode-name)
  "Name of the go test.  MODE-NAME is unused."
  "*Go test*")

(defun gotn--run-test-as-compilation (cmd)
  "Run CMD in gotn-mode."
  (compilation-start cmd
                     'gotn-mode
                     'gotn--compilation-name))

;;;###autoload
(defun gotn-run-test-package ()
  "Run go test of current package."
  (interactive "d")
  (gotn--run-test-as-compilation go-test-package-command))

;;;###autoload
(defun gotn-run-test (point)
  "Run go test at POINT."
  (interactive "d")
  (condition-case nil
      (let ((gotn-out (gotn--call point)))
        (if (= (car gotn-out) 0)
            (gotn--run-test-as-compilation (concat go-test-case-command " ^" (car (cdr gotn-out)) "$"))
          (if go-test-package-test-fallback
              (gotn--run-test-as-compilation go-test-package-command)
            (message (format "Could not run gotn binary: %s" (cdr gotn-out))))))))

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

(defvar gotn-mode-map
  (nconc (make-sparse-keymap) compilation-mode-map)
  "Keymap for gotn major mode.")

(define-derived-mode gotn-mode compilation-mode "gotn"
  "Major mode for the gotn compilation buffer."
  (use-local-map gotn-mode-map)
  (setq major-mode 'gotn-mode)
  (setq mode-name "gotn")
  (setq-local truncate-lines t)
  (font-lock-add-keywords nil nil))

(provide 'gotn)

;;; gotn.el ends here
