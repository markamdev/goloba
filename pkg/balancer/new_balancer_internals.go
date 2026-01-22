package balancer

import (
	"context"
	"net"

	"github.com/markamdev/goloba/pkg/logger"
	"github.com/pkg/errors"
)

const (
	defaultMaxLoopCount = 100
)

type balancerImpl struct {
	cfg       Config
	logger    logger.GLBLogger
	listener  net.Listener
	loopCount int // to avoid infinite loops in case of repetitive error
	ctxCancel context.CancelFunc
}

func (b *balancerImpl) startListener() error {
	b.logger.Debug("starting listener", "port", b.cfg.Port)

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"), // TODO consider binding address selection
		Port: int(b.cfg.Port),
	})
	if err != nil {
		return errors.Wrap(err, "failed to start TCP listener")
	}
	b.listener = listener
	b.loopCount = 0

	var connCtx context.Context
	connCtx, b.ctxCancel = context.WithCancel(context.Background())

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				b.logger.Debug("listener closed")
				return nil
			}
			b.loopCount++
			if b.loopCount > defaultMaxLoopCount {
				b.logger.Error("too many consecutive errors accepting connections, stopping listener")
				return listener.Close()
			}
			b.logger.Error("failed to accept connection", "error", err.Error())
			continue
		}
		b.loopCount = 0
		go b.handleConnection(connCtx, conn)
	}
}

// TODO implement connection handling and closing when listener closed
func (b *balancerImpl) handleConnection(ctx context.Context, conn net.Conn) {
	b.logger.Debug("accepted new connection", "remote_addr", conn.RemoteAddr().String())
	defer conn.Close()
}

func (b *balancerImpl) stopListener() error {
	if b.listener != nil {
		b.logger.Debug("closing listener")
		lstCopy := b.listener
		b.listener = nil
		return lstCopy.Close()
	}
	return nil
}
