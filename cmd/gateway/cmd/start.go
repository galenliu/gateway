package cmd

import (
	"context"
	"fmt"
	"github.com/galenliu/gateway"
	"github.com/galenliu/gateway/pkg/constant"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	serviceName = "gatewayService"
)

func (c *command) initStartCmd() (err error) {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start WebThings gateway",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) > 0 {
				return cmd.Help()
			}

			v := strings.ToLower(c.config.GetString(optionNameVerbosity))
			logger, err := newLogger(cmd, v)
			if err != nil {
				return fmt.Errorf("new logger: %v", err)
			}

			isWindowsService, err := isWindowsService()
			if err != nil {
				return fmt.Errorf("failed to determine if we are running in service: %w", err)
			}

			if isWindowsService {
				var err error
				logger, err = createWindowsEventLogger(serviceName, logger)
				if err != nil {
					return fmt.Errorf("failed to create windows logger %w", err)
				}
			}

			logger.Infof("gateway version %v", constant.Version)

			ctx, cancelFunc := context.WithCancel(context.Background())
			gw, err := gateway.NewGateway(ctx, gateway.Config{
				BaseDir:          c.config.GetString(optionNameDataDir),
				AttachAddonsDir:  c.config.GetString(optionNameAttachAddonDirs),
				RemoveBeforeOpen: c.config.GetBool(optionNameDBRemoveBeforeOpen),
				Verbosity:        c.config.GetString(optionNameVerbosity),
				AddonUrls:        c.config.GetStringSlice(optionNameAddonUrls),
				IPCPort:          ":" + strconv.Itoa(c.config.GetInt(optionNameIpcPort)),
				RPCPort:          ":" + strconv.Itoa(c.config.GetInt(optionNameRpcPort)),
				HttpAddr:         ":" + strconv.Itoa(c.config.GetInt(optionNameHttpPort)),
				HttpsAddr:        ":" + strconv.Itoa(c.config.GetInt(optionNameHttpsPort)),
				LogRotateDays:    c.config.GetInt(optionLogRotateDays),
				HomeKitPin:       c.config.GetString(optionHomeKitPin),
				HomeKitEnable:    c.config.GetBool(optionHomeKitEnable),
			}, logger)
			if err != nil {
				cancelFunc()
				return err
			}

			// Wait for termination or interrupt signals.
			// We want to clean up things at the end.
			interruptChannel := make(chan os.Signal, 1)
			signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM)

			p := &program{
				start: func() {
					// Block main goroutine until it is interrupted
					sig := <-interruptChannel
					logger.Debugf("received signal: %v", sig)
					logger.Info("shutting down")
					return
				},
				stop: func() {
					// Shutdown

					go func() {

						ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
						defer cancel()
						if err := gw.Shutdown(ctx); err != nil {
							logger.Errorf("shutdown: %v", err)
						}
					}()

					// If shutdown function is blocking too long,
					// allow process termination by receiving another signal.
					select {
					case sig := <-interruptChannel:
						cancelFunc()
						logger.Debugf("received signal: %v", sig)
						return
					case <-ctx.Done():
						cancelFunc()
					}
				},
			}

			if isWindowsService {
				s, err := service.New(p, &service.Config{
					Name:        serviceName,
					DisplayName: "Gateway",
					Description: "WebThings, gateway service.",
				})
				if err != nil {
					return err
				}
				if err = s.Run(); err != nil {
					return err
				}
			} else {
				// start blocks until some interrupt is received
				p.start()
				p.stop()
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.config.BindPFlags(cmd.Flags())
		},
	}

	c.setAllFlags(cmd)
	c.root.AddCommand(cmd)
	return nil
}

type program struct {
	start func()
	stop  func()
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.start()
	return nil
}

func (p *program) Stop(s service.Service) error {
	p.stop()
	return nil
}
