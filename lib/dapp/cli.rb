require 'mixlib/cli'

module Dapp
  # CLI
  class CLI
    include Mixlib::CLI
    include CliHelper

    SUBCOMMANDS = %w(build smartpush push list show flush).freeze

    banner <<BANNER.freeze
Usage: dapp [options] sub-command [sub-command options]

Available subcommands: (for details, dapp SUB-COMMAND --help)

dapp build [options] [PATTERN ...]
dapp push [options] [PATTERN] REPO
dapp smartpush [options] [PATTERN ...] REPOPREFIX
dapp list [options] [PATTERN ...]
dapp show [options] [PATTERN ...]
dapp flush

Options:
BANNER

    option :version,
           long: '--version',
           description: 'Show version',
           on: :tail,
           boolean: true,
           proc: proc { puts "dapp: #{Dapp::VERSION}" },
           exit: 0

    option :help,
           short: '-h',
           long: '--help',
           description: 'Show this message',
           on: :tail,
           boolean: true,
           show_options: true,
           exit: 0

    def initialize(*args)
      super(*args)

      opt_parser.program_name = 'dapp'
      opt_parser.version = Dapp::VERSION
    end

    def run(argv = ARGV)
      argv, subcommand, subcommand_argv = self.class.parse_subcommand(self, argv)
      self.class.parse_options(self, argv)
      self.class.run_subcommand self, subcommand, subcommand_argv
    end
  end
end
