package server

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Spexso/CSE343-Online-Library-System/backend/libware/errlist"
	"github.com/Spexso/CSE343-Online-Library-System/backend/libware/helpers"
	"github.com/Spexso/CSE343-Online-Library-System/backend/libware/server/requests"
	"github.com/Spexso/CSE343-Online-Library-System/backend/libware/server/responses"
)

func (l *LibraryHandler) userHandler() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/isbn-profile", l.isbnProfile)
	router.HandleFunc("/user-profile", l.userProfile)
	router.HandleFunc("/change-user-name", l.changeUserName)
	router.HandleFunc("/change-user-email", l.changeUserEmail)
	router.HandleFunc("/change-user-phone", l.changeUserPhone)
	router.HandleFunc("/change-user-password", l.changeUserPassword)
	router.HandleFunc("/enqueue", l.Enqueue)
	router.HandleFunc("/dequeue", l.Dequeue)
	router.HandleFunc("/suspended-until", l.SuspendedUntil)
	router.HandleFunc("/queued-books", l.QueuedBooks)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subject, err := l.authorize(w, r, l.userSecret)
		if err != nil {
			log.Printf("error: user: %v", err)
			return
		}

		r.Header.Add("Subject", subject)
		router.ServeHTTP(w, r)
	})
}

func (l *LibraryHandler) userProfile(w http.ResponseWriter, r *http.Request) {
	idString := r.Header.Get("Subject")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: user-profile: %v", err)
		return
	}

	name, surname, email, phone, err := l.db.UserProfile(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: user-profile: %v", err)
		return
	}

	response := responses.UserProfile{
		Name:    name,
		Surname: surname,
		Email:   email,
		Phone:   phone,
	}
	helpers.WriteResponse(w, response)

	w.WriteHeader(http.StatusOK)
}

func (l *LibraryHandler) changeUserName(w http.ResponseWriter, r *http.Request) {
	var req requests.ChangeUserName
	err := helpers.ReadRequest(r.Body, &req)
	if err != nil || req.NewName == "" || req.NewSurname == "" {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrJsonDecoder)
		if err == nil {
			err = errlist.ErrJsonDecoder
		}
		log.Printf("error: change-user-name: %v", err)
		return
	}

	idString := r.Header.Get("Subject")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: change-user-name: %v", err)
		return
	}

	err = l.db.ChangeUserName(id, req.NewName, req.NewSurname)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: change-user-name: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (l *LibraryHandler) changeUserEmail(w http.ResponseWriter, r *http.Request) {
	var req requests.ChangeUserEmail
	err := helpers.ReadRequest(r.Body, &req)
	if err != nil || req.Password == "" || req.NewEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrJsonDecoder)
		if err == nil {
			err = errlist.ErrJsonDecoder
		}
		log.Printf("error: change-user-email: %v", err)
		return
	}

	idString := r.Header.Get("Subject")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: change-user-email: %v", err)
		return
	}

	err = l.db.ChangeUserEmail(id, req.Password, req.NewEmail)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errors.Is(err, errlist.ErrEmailExist) {
			helpers.WriteError(w, errlist.ErrEmailExist)
		} else if errors.Is(err, errlist.ErrInvalidPassword) {
			helpers.WriteError(w, errlist.ErrInvalidPassword)
		} else {
			helpers.WriteError(w, errlist.ErrGeneric)
		}
		log.Printf("error: change-user-email: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (l *LibraryHandler) changeUserPhone(w http.ResponseWriter, r *http.Request) {
	var req requests.ChangeUserPhone
	err := helpers.ReadRequest(r.Body, &req)
	if err != nil || req.Password == "" || req.NewPhone == "" {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrJsonDecoder)
		if err == nil {
			err = errlist.ErrJsonDecoder
		}
		log.Printf("error: change-user-phone: %v", err)
		return
	}

	idString := r.Header.Get("Subject")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: change-user-phone: %v", err)
		return
	}

	err = l.db.ChangeUserPhone(id, req.Password, req.NewPhone)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errors.Is(err, errlist.ErrInvalidPassword) {
			helpers.WriteError(w, errlist.ErrInvalidPassword)
		} else {
			helpers.WriteError(w, errlist.ErrGeneric)
		}
		log.Printf("error: change-user-phone: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (l *LibraryHandler) changeUserPassword(w http.ResponseWriter, r *http.Request) {
	var req requests.ChangeUserPassword
	err := helpers.ReadRequest(r.Body, &req)
	if err != nil || req.OldPassword == "" || req.NewPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrJsonDecoder)
		if err == nil {
			err = errlist.ErrJsonDecoder
		}
		log.Printf("error: change-user-password: %v", err)
		return
	}

	idString := r.Header.Get("Subject")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: change-user-password: %v", err)
		return
	}

	err = l.db.ChangeUserPassword(id, req.OldPassword, req.NewPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errors.Is(err, errlist.ErrInvalidPassword) {
			helpers.WriteError(w, errlist.ErrInvalidPassword)
		} else {
			helpers.WriteError(w, errlist.ErrGeneric)
		}
		log.Printf("error: change-user-password: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (l *LibraryHandler) Enqueue(w http.ResponseWriter, r *http.Request) {
	var req requests.Enqueue
	err := helpers.ReadRequest(r.Body, &req)
	if err != nil || req.Isbn == "" {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrJsonDecoder)
		if err == nil {
			err = errlist.ErrJsonDecoder
		}
		log.Printf("error: user-enqueue: %v", err)
		return
	}

	idString := r.Header.Get("Subject")

	userId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: user-enqueue: %v", err)
		return
	}

	err = l.db.UserEnqueue(userId, req.Isbn)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errors.Is(err, errlist.ErrUserSuspended) {
			helpers.WriteError(w, errlist.ErrUserSuspended)
		} else if errors.Is(err, errlist.ErrPastDue) {
			helpers.WriteError(w, errlist.ErrPastDue)
		} else if errors.Is(err, errlist.ErrUserInQueue) {
			helpers.WriteError(w, errlist.ErrUserInQueue)
		} else if errors.Is(err, errlist.ErrIsbnNotExist) {
			helpers.WriteError(w, errlist.ErrIsbnNotExist)
		} else if errors.Is(err, errlist.ErrAlreadyBorrowed) {
			helpers.WriteError(w, errlist.ErrAlreadyBorrowed)
		} else {
			helpers.WriteError(w, errlist.ErrGeneric)
		}
		log.Printf("error: user-enqueue: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (l *LibraryHandler) Dequeue(w http.ResponseWriter, r *http.Request) {
	var req requests.Dequeue
	err := helpers.ReadRequest(r.Body, &req)
	if err != nil || req.Isbn == "" {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrJsonDecoder)
		if err == nil {
			err = errlist.ErrJsonDecoder
		}
		log.Printf("error: user-dequeue: %v", err)
		return
	}

	idString := r.Header.Get("Subject")

	userId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: user-dequeue: %v", err)
		return
	}

	err = l.db.UserDequeue(userId, req.Isbn)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errors.Is(err, errlist.ErrUserNotInQueue) {
			helpers.WriteError(w, errlist.ErrUserNotInQueue)
		} else if errors.Is(err, errlist.ErrIsbnNotExist) {
			helpers.WriteError(w, errlist.ErrIsbnNotExist)
		} else {
			helpers.WriteError(w, errlist.ErrGeneric)
		}
		log.Printf("error: user-dequeue: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (l *LibraryHandler) SuspendedUntil(w http.ResponseWriter, r *http.Request) {
	idString := r.Header.Get("Subject")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: suspended-until: %v", err)
		return
	}

	timestamp, err := l.db.UserSuspendedUntil(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: suspended-until: %v", err)
		return
	}

	response := responses.SuspendedUntil{
		Timestamp: strconv.FormatInt(timestamp, 10),
	}
	helpers.WriteResponse(w, response)

	w.WriteHeader(http.StatusOK)
}

func (l *LibraryHandler) QueuedBooks(w http.ResponseWriter, r *http.Request) {
	idString := r.Header.Get("Subject")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: queued-books: %v", err)
		return
	}

	isbns, err := l.db.UserQueuedBooks(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: queued-books: %v", err)
		return
	}
	for _, isbn := range isbns {
		l.db.IsbnQueueEnforceInvariants(isbn)
	}

	isbns, availableBooks, indexes, validUntils, err := l.db.UserQueuedBooksDetailed(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		helpers.WriteError(w, errlist.ErrGeneric)
		log.Printf("error: queued-books: %v", err)
		return
	}

	response := responses.QueuedBooks{
		Entries: make([]responses.QueuedBookEntry, len(isbns)),
	}
	for i := range response.Entries {
		response.Entries[i].Isbn = isbns[i]
		response.Entries[i].AvailableBooks = strconv.FormatInt(availableBooks[i], 10)
		response.Entries[i].Position = strconv.FormatInt(indexes[i]+1, 10)
		response.Entries[i].ValidUntil = strconv.FormatInt(validUntils[i], 10)
	}

	helpers.WriteResponse(w, response)

	w.WriteHeader(http.StatusOK)
}
