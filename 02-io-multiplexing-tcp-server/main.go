package main

import (
	"fmt"
	"net"
	"syscall"
)

func main() {
	port := net.TCPAddr{Port: 8080}
	ln, err := net.ListenTCP("tcp", &port)
	if err != nil {
		fmt.Println("listener error", err)
		return
	}
	fmt.Println("server started on port ", port.Port)

	listenerFile, err := ln.File()
	if err != nil {
		fmt.Println("fail to get listener fd", err)
	}
	defer listenerFile.Close()

	serverFd := int(listenerFile.Fd())
	syscall.SetNonblock(serverFd, true)

	// Tạo 1 Epoll Instance để nhờ OS theo dõi các socket
	// epfd nhận được chính là File Descriptor của bản thân cây epoll đó.
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		fmt.Println("fail to create Epoll")
	}

	//Bây giờ cần bỏ fd của listener vào epoll để nó theo dõi sự kiện Có Client kết nối tới (EPOLLIN).

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN
	event.Fd = int32(serverFd)

	// Thêm vào epoll bằng syscall.EpollCtl:
	syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, serverFd, &event)

	// ## Dựng Vòng lặp Event Loop
	// tạo sẵn 1 mảng để OS chép các sự kiện sẵn sàng vào:
	events := make([]syscall.EpollEvent, 100) // chứa tối đa 100 sự kiện mỗi lần chờ

	for {
		n, err := syscall.EpollWait(epfd, events, -1)
		if err != nil {
			// Trong syscall có thể gặp lỗi EINTR (bị ngắt bởi system signal), có thể continue
			continue
		}

		for i := 0; i < n; i++ {
			ev := events[i]

			if ev.Fd == int32(serverFd) {
				newFd, _, err := syscall.Accept(int(ev.Fd))
				if err != nil {
					fmt.Println("fail to accept new socket")
					continue
				}
				syscall.SetNonblock(newFd, true)
				fmt.Printf("new client connected on Fd: %v\n", newFd)

				var event syscall.EpollEvent
				event.Events = syscall.EPOLLIN
				event.Fd = int32(newFd)

				// Thêm vào epoll bằng syscall.EpollCtl:
				syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, newFd, &event)
			} else {
				buffer := make([]byte, 1024)
				nBytes, err := syscall.Read(int(ev.Fd), buffer)

				// Nếu có lỗi HOẶC client đóng kết nối (nBytes == 0)
				if err != nil || nBytes == 0 {
					if err != nil {
						fmt.Println("read error on fd", ev.Fd, ":", err)
					} else {
						fmt.Println("client disconnected", ev.Fd)
					}

					// Luôn dọn dẹp socket kẹt
					syscall.EpollCtl(epfd, syscall.EPOLL_CTL_DEL, int(ev.Fd), nil)
					syscall.Close(int(ev.Fd))
					continue

				}

				data := string(buffer[:nBytes])
				fmt.Print("Data: ", data)

			}

		}
	}

}
