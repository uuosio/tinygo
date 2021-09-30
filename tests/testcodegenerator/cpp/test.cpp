#include <eosio/eosio.hpp>
#include <eosio/crypto.hpp>
#include <eosio/binary_extension.hpp>

using namespace std;
struct test_struct {
    string a;
    eosio::binary_extension<eosio::checksum256> b;
    eosio::binary_extension<eosio::checksum256> c;
};

struct test_optional {
    string a;
    std::optional<eosio::checksum256> b;
    std::optional<eosio::checksum256> c;
};

struct test_optional_plus_extension {
    string a;
    std::optional<eosio::checksum256> c;
    eosio::binary_extension<eosio::checksum256> b;
};

extern "C" void apply(uint64_t receiver, uint64_t first_receiver, uint64_t action) {
    if (action == "testext"_n.value) {
        auto args = eosio::unpack_action_data<test_struct>();
        eosio::check(args.b.has_value(), "no args");
        auto& b = args.b.value();
        eosio::print("value: ", b, "\n");
        auto& c = args.c.value();
        eosio::print("value: ", c, "\n");
    } else if (action == "testext2"_n.value) {
        auto args = eosio::unpack_action_data<test_struct>();
        eosio::check(!args.b.has_value(), "!args.b.has_value()");
        eosio::check(!args.c.has_value(), "!args.c.has_value()");
    } else if (action == "testopt"_n.value) {
        auto args = eosio::unpack_action_data<test_optional>();
        eosio::check(args.b.has_value(), "args.b.is_valid()");
        eosio::check(args.c.has_value(), "args.c.is_valid()");
        eosio::print(args.b.value());
        eosio::print(args.c.value());
    } else if (action == "testopt2"_n.value) {
        auto args = eosio::unpack_action_data<test_optional>();
        eosio::check(!args.b.has_value(), "!args.b.has_value()");
        eosio::check(!args.c.has_value(), "!args.c.has_value()");
    } else if (action == "testcombine"_n.value) {
        auto args = eosio::unpack_action_data<test_optional_plus_extension>();
        eosio::check(args.b.has_value(), "args.b.has_value()");
        eosio::check(args.c.has_value(), "args.c.has_value()");
    }
}
